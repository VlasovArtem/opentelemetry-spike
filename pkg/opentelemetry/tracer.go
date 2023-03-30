package opentelemetry

import (
	"context"
	crand "crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.18.0"
	"go.opentelemetry.io/otel/trace"
	"math/rand"
	"os"
	"sync"
)

const traceHeaderName = "X-Correlation-ID"

type TraceHeaderContext struct{}

func (t TraceHeaderContext) Inject(ctx context.Context, carrier propagation.TextMapCarrier) {
	value := ctx.Value(traceHeaderName)
	if header, ok := value.(string); ok {
		carrier.Set(traceHeaderName, header)
	}
}

func (t TraceHeaderContext) Extract(ctx context.Context, carrier propagation.TextMapCarrier) context.Context {
	header := carrier.Get(traceHeaderName)
	if header != "" {
		return context.WithValue(ctx, traceHeaderName, header)
	}
	return ctx
}

func (t TraceHeaderContext) Fields() []string {
	return []string{traceHeaderName}
}

type customIDGenerator struct {
	sync.Mutex
	randSource *rand.Rand
}

// NewSpanID returns a non-zero span ID from a randomly-chosen sequence.
func (gen *customIDGenerator) NewSpanID(ctx context.Context, traceID trace.TraceID) trace.SpanID {
	gen.Lock()
	defer gen.Unlock()
	sid := trace.SpanID{}
	_, _ = gen.randSource.Read(sid[:])
	return sid
}

// NewIDs returns a non-zero trace ID and a non-zero span ID from a
// randomly-chosen sequence.
func (gen *customIDGenerator) NewIDs(ctx context.Context) (trace.TraceID, trace.SpanID) {
	gen.Lock()
	defer gen.Unlock()
	tid := trace.TraceID{}
	_, _ = gen.randSource.Read(tid[:])
	value := ctx.Value(traceHeaderName)
	if header, ok := value.(string); ok {
		hexRequestId := hex.EncodeToString([]byte(header))
		tidFromHeader, err := trace.TraceIDFromHex(hexRequestId)
		if err == nil {
			tid = tidFromHeader
		} else {
			log.Error().Err(err).Msg("failed to parse trace id from header")
		}
	}
	sid := trace.SpanID{}
	_, _ = gen.randSource.Read(sid[:])
	return tid, sid
}

func defaultIDGenerator() sdktrace.IDGenerator {
	gen := &customIDGenerator{}
	var rngSeed int64
	_ = binary.Read(crand.Reader, binary.LittleEndian, &rngSeed)
	gen.randSource = rand.New(rand.NewSource(rngSeed))
	return gen
}

func EnableGlobalFileTracer(serviceName string) func(ctx context.Context) error {
	f, err := os.Create(serviceName + ".txt")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create file")
	}

	exporter, err := NewWriterExporter(f)

	if err != nil {
		log.Fatal().Err(err).Msg("failed to create exporter")
	}

	cleanup := initTracer(serviceName, exporter)

	return func(ctx context.Context) error {
		f.Close()
		return cleanup(ctx)
	}
}

func EnabledGlobalGrpcTracer(serviceName string, url string, insecure bool) func(ctx context.Context) error {
	exporter, err := NewGrpcTraceExporter(url, insecure)

	if err != nil {
		log.Fatal().Err(err).Msg("failed to create exporter")
	}

	cleanup := initTracer(serviceName, exporter)

	return cleanup
}

func initTracer(serviceName string, exporter sdktrace.SpanExporter) func(ctx context.Context) error {
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(newResource(serviceName)),
		sdktrace.WithIDGenerator(defaultIDGenerator()),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			TraceHeaderContext{},
		),
	)

	return tp.Shutdown
}

func newResource(serviceName string) *resource.Resource {
	return resource.NewSchemaless(
		semconv.ServiceName(serviceName),
		semconv.ServiceVersion("v0.1.0"),
		attribute.String("environment", "demo"),
	)
}
