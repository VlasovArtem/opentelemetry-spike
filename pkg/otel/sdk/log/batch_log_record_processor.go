package log

import (
	"context"
	"runtime"
	"spike-go-opentelemetry-logging/pkg/otel/internal/global"
	"spike-go-opentelemetry-logging/pkg/otel/log"
	"spike-go-opentelemetry-logging/pkg/otel/sdk/internal/env"
	"sync"
	"sync/atomic"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

// Defaults for BatchLogRecordProcessorOptions.
const (
	DefaultMaxQueueSize       = 2048
	DefaultScheduleDelay      = 5000
	DefaultExportTimeout      = 30000
	DefaultMaxExportBatchSize = 512
)

// BatchLogRecordProcessorOption configures a BatchSpanProcessor.
type BatchLogRecordProcessorOption func(o *BatchLogRecordProcessorOptions)

// BatchLogRecordProcessorOptions is configuration settings for a
// BatchSpanProcessor.
type BatchLogRecordProcessorOptions struct {
	MaxQueueSize int

	BatchTimeout time.Duration

	ExportTimeout time.Duration

	MaxExportBatchSize int

	BlockOnQueueFull bool
}

type batchLogRecordProcessor struct {
	e LogExporter
	o BatchLogRecordProcessorOptions

	queue   chan log.ReadableLogRecord
	dropped uint32

	batch      []log.ReadableLogRecord
	batchMutex sync.Mutex
	timer      *time.Timer
	stopWait   sync.WaitGroup
	stopOnce   sync.Once
	stopCh     chan struct{}
}

func NewBatchLogProcessorProcessor(exporter LogExporter, options ...BatchLogRecordProcessorOption) LogRecordProcessor {
	maxQueueSize := env.BatchLogRecordProcessorMaxQueueSize(DefaultMaxQueueSize)
	maxExportBatchSize := env.BatchLogRecordProcessorMaxExportBatchSize(DefaultMaxExportBatchSize)

	if maxExportBatchSize > maxQueueSize {
		if DefaultMaxExportBatchSize > maxQueueSize {
			maxExportBatchSize = maxQueueSize
		} else {
			maxExportBatchSize = DefaultMaxExportBatchSize
		}
	}

	o := BatchLogRecordProcessorOptions{
		BatchTimeout:       time.Duration(env.BatchLogRecordProcessorScheduleDelay(DefaultScheduleDelay)) * time.Millisecond,
		ExportTimeout:      time.Duration(env.BatchLogRecordProcessorExportTimeout(DefaultExportTimeout)) * time.Millisecond,
		MaxQueueSize:       maxQueueSize,
		MaxExportBatchSize: maxExportBatchSize,
	}
	for _, opt := range options {
		opt(&o)
	}
	blrp := &batchLogRecordProcessor{
		e:      exporter,
		o:      o,
		batch:  make([]log.ReadableLogRecord, 0, o.MaxExportBatchSize),
		timer:  time.NewTimer(o.BatchTimeout),
		queue:  make(chan log.ReadableLogRecord, o.MaxQueueSize),
		stopCh: make(chan struct{}),
	}

	blrp.stopWait.Add(1)
	go func() {
		defer blrp.stopWait.Done()
		blrp.processQueue()
		blrp.drainQueue()
	}()

	return blrp
}

// OnEmit method enqueues a ReadOnlySpan for later processing.
func (bsp *batchLogRecordProcessor) OnEmit(s log.ReadableLogRecord) {
	// Do not enqueue spans if we are just going to drop them.
	if bsp.e == nil {
		return
	}
	bsp.enqueue(s)
}

// Shutdown flushes the queue and waits until all spans are processed.
// It only executes once. Subsequent call does nothing.
func (bsp *batchLogRecordProcessor) Shutdown(ctx context.Context) error {
	var err error
	bsp.stopOnce.Do(func() {
		wait := make(chan struct{})
		go func() {
			close(bsp.stopCh)
			bsp.stopWait.Wait()
			if bsp.e != nil {
				if err := bsp.e.Shutdown(ctx); err != nil {
					otel.Handle(err)
				}
			}
			close(wait)
		}()
		// Wait until the wait group is done or the context is cancelled
		select {
		case <-wait:
		case <-ctx.Done():
			err = ctx.Err()
		}
	})
	return err
}

type forceFlushSpan struct {
	log.ReadableLogRecord
	flushed chan struct{}
}

func (f forceFlushSpan) SpanContext() trace.SpanContext {
	return trace.NewSpanContext(trace.SpanContextConfig{TraceFlags: trace.FlagsSampled})
}

// ForceFlush exports all ended spans that have not yet been exported.
func (bsp *batchLogRecordProcessor) ForceFlush(ctx context.Context) error {
	var err error
	if bsp.e != nil {
		flushCh := make(chan struct{})
		if bsp.enqueueBlockOnQueueFull(ctx, forceFlushSpan{flushed: flushCh}) {
			select {
			case <-flushCh:
				// Processed any items in queue prior to ForceFlush being called
			case <-ctx.Done():
				return ctx.Err()
			}
		}

		wait := make(chan error)
		go func() {
			wait <- bsp.exportLogRecords(ctx)
			close(wait)
		}()
		// Wait until the export is finished or the context is cancelled/timed out
		select {
		case err = <-wait:
		case <-ctx.Done():
			err = ctx.Err()
		}
	}
	return err
}

// WithMaxQueueSize returns a BatchLogRecordProcessorOption that configures the
// maximum queue size allowed for a BatchSpanProcessor.
func WithMaxQueueSize(size int) BatchLogRecordProcessorOption {
	return func(o *BatchLogRecordProcessorOptions) {
		o.MaxQueueSize = size
	}
}

// WithMaxExportBatchSize returns a BatchLogRecordProcessorOption that configures
// the maximum export batch size allowed for a BatchSpanProcessor.
func WithMaxExportBatchSize(size int) BatchLogRecordProcessorOption {
	return func(o *BatchLogRecordProcessorOptions) {
		o.MaxExportBatchSize = size
	}
}

// WithBatchTimeout returns a BatchLogRecordProcessorOption that configures the
// maximum delay allowed for a BatchSpanProcessor before it will export any
// held span (whether the queue is full or not).
func WithBatchTimeout(delay time.Duration) BatchLogRecordProcessorOption {
	return func(o *BatchLogRecordProcessorOptions) {
		o.BatchTimeout = delay
	}
}

// WithExportTimeout returns a BatchLogRecordProcessorOption that configures the
// amount of time a BatchSpanProcessor waits for an exporter to export before
// abandoning the export.
func WithExportTimeout(timeout time.Duration) BatchLogRecordProcessorOption {
	return func(o *BatchLogRecordProcessorOptions) {
		o.ExportTimeout = timeout
	}
}

// WithBlocking returns a BatchLogRecordProcessorOption that configures a
// BatchSpanProcessor to wait for enqueue operations to succeed instead of
// dropping data when the queue is full.
func WithBlocking() BatchLogRecordProcessorOption {
	return func(o *BatchLogRecordProcessorOptions) {
		o.BlockOnQueueFull = true
	}
}

// exportLogRecords is a subroutine of processing and draining the queue.
func (bsp *batchLogRecordProcessor) exportLogRecords(ctx context.Context) error {
	bsp.timer.Reset(bsp.o.BatchTimeout)

	bsp.batchMutex.Lock()
	defer bsp.batchMutex.Unlock()

	if bsp.o.ExportTimeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, bsp.o.ExportTimeout)
		defer cancel()
	}

	if l := len(bsp.batch); l > 0 {
		global.Debug("exporting spans", "count", len(bsp.batch), "total_dropped", atomic.LoadUint32(&bsp.dropped))
		err := bsp.e.ExportLogRecords(ctx, bsp.batch)

		// A new batch is always created after exporting, even if the batch failed to be exported.
		//
		// It is up to the exporter to implement any type of retry logic if a batch is failing
		// to be exported, since it is specific to the protocol and backend being sent to.
		bsp.batch = bsp.batch[:0]

		if err != nil {
			return err
		}
	}
	return nil
}

// processQueue removes spans from the `queue` channel until processor
// is shut down. It calls the exporter in batches of up to MaxExportBatchSize
// waiting up to BatchTimeout to form a batch.
func (bsp *batchLogRecordProcessor) processQueue() {
	defer bsp.timer.Stop()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for {
		select {
		case <-bsp.stopCh:
			return
		case <-bsp.timer.C:
			if err := bsp.exportLogRecords(ctx); err != nil {
				otel.Handle(err)
			}
		case sd := <-bsp.queue:
			if ffs, ok := sd.(forceFlushSpan); ok {
				close(ffs.flushed)
				continue
			}
			bsp.batchMutex.Lock()
			bsp.batch = append(bsp.batch, sd)
			shouldExport := len(bsp.batch) >= bsp.o.MaxExportBatchSize
			bsp.batchMutex.Unlock()
			if shouldExport {
				if !bsp.timer.Stop() {
					<-bsp.timer.C
				}
				if err := bsp.exportLogRecords(ctx); err != nil {
					otel.Handle(err)
				}
			}
		}
	}
}

// drainQueue awaits the any caller that had added to bsp.stopWait
// to finish the enqueue, then exports the final batch.
func (bsp *batchLogRecordProcessor) drainQueue() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for {
		select {
		case sd := <-bsp.queue:
			if sd == nil {
				if err := bsp.exportLogRecords(ctx); err != nil {
					otel.Handle(err)
				}
				return
			}

			bsp.batchMutex.Lock()
			bsp.batch = append(bsp.batch, sd)
			shouldExport := len(bsp.batch) == bsp.o.MaxExportBatchSize
			bsp.batchMutex.Unlock()

			if shouldExport {
				if err := bsp.exportLogRecords(ctx); err != nil {
					otel.Handle(err)
				}
			}
		default:
			close(bsp.queue)
		}
	}
}

func (bsp *batchLogRecordProcessor) enqueue(sd log.ReadableLogRecord) {
	ctx := context.TODO()
	if bsp.o.BlockOnQueueFull {
		bsp.enqueueBlockOnQueueFull(ctx, sd)
	} else {
		bsp.enqueueDrop(ctx, sd)
	}
}

func recoverSendOnClosedChan() {
	x := recover()
	switch err := x.(type) {
	case nil:
		return
	case runtime.Error:
		if err.Error() == "send on closed channel" {
			return
		}
	}
	panic(x)
}

func (bsp *batchLogRecordProcessor) enqueueBlockOnQueueFull(ctx context.Context, sd log.ReadableLogRecord) bool {
	defer recoverSendOnClosedChan()

	select {
	case <-bsp.stopCh:
		return false
	default:
	}

	select {
	case bsp.queue <- sd:
		return true
	case <-ctx.Done():
		return false
	}
}

func (bsp *batchLogRecordProcessor) enqueueDrop(ctx context.Context, sd log.ReadableLogRecord) bool {
	defer recoverSendOnClosedChan()

	select {
	case <-bsp.stopCh:
		return false
	default:
	}

	select {
	case bsp.queue <- sd:
		return true
	default:
		atomic.AddUint32(&bsp.dropped, 1)
	}
	return false
}

// MarshalLog is the marshaling function used by the logging system to represent this exporter.
func (bsp *batchLogRecordProcessor) MarshalLog() interface{} {
	return struct {
		Type        string
		LogExporter LogExporter
		Config      BatchLogRecordProcessorOptions
	}{
		Type:        "BatchSpanProcessor",
		LogExporter: bsp.e,
		Config:      bsp.o,
	}
}
