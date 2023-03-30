package log

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/resource"
	"spike-go-opentelemetry-logging/pkg/otel/exporters/otlp/otlplog/global"
	"spike-go-opentelemetry-logging/pkg/otel/log"
	"sync"
	"sync/atomic"
)

const (
	defaultLoggerName = "go.opentelemetry.io/otel/sdk/log"
)

type LoggerProvider struct {
	mu             sync.Mutex
	namedTracer    map[instrumentation.Scope]*logger
	logsProcessors atomic.Value
	isShutdown     bool
	resource       *resource.Resource
}

func NewLoggerProvider(opts ...LoggerProviderOption) *LoggerProvider {
	o := loggerProviderConfig{}

	for _, opt := range opts {
		o = opt.apply(o)
	}

	o = ensureValidLoggerProviderConfig(o)

	tp := &LoggerProvider{
		namedTracer: make(map[instrumentation.Scope]*logger),
		resource:    o.resource,
	}
	global.Info("TracerProvider created", "config", o)

	spss := logRecordProcessorStates{}
	for _, sp := range o.processors {
		spss = append(spss, newLogRecordProcessorState(sp))
	}
	tp.logsProcessors.Store(spss)

	return tp
}

func (p *LoggerProvider) Shutdown(ctx context.Context) error {
	spss := p.logsProcessors.Load().(logRecordProcessorStates)
	if len(spss) == 0 {
		return nil
	}

	p.mu.Lock()
	defer p.mu.Unlock()
	p.isShutdown = true

	var retErr error
	for _, sps := range spss {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		var err error
		sps.state.Do(func() {
			err = sps.sp.Shutdown(ctx)
		})
		if err != nil {
			if retErr == nil {
				retErr = err
			} else {
				// Poor man's list of errors
				retErr = fmt.Errorf("%v; %v", retErr, err)
			}
		}
	}
	p.logsProcessors.Store(logRecordProcessorStates{})
	return retErr
}

type loggerProviderConfig struct {
	processors []LogRecordProcessor

	// resource contains attributes representing an entity that produces telemetry.
	resource *resource.Resource
}

type LoggerProviderOption interface {
	apply(loggerProviderConfig) loggerProviderConfig
}

type loggerProviderOptionFunc func(loggerProviderConfig) loggerProviderConfig

func (fn loggerProviderOptionFunc) apply(cfg loggerProviderConfig) loggerProviderConfig {
	return fn(cfg)
}

func ensureValidLoggerProviderConfig(cfg loggerProviderConfig) loggerProviderConfig {
	if cfg.resource == nil {
		cfg.resource = resource.Default()
	}
	return cfg
}

func (p *LoggerProvider) Logger(name string, opts ...log.LoggerOption) log.Logger {
	c := log.NewLoggerConfig(opts...)

	p.mu.Lock()
	defer p.mu.Unlock()
	if name == "" {
		name = defaultLoggerName
	}
	is := instrumentation.Scope{
		Name:      name,
		Version:   c.InstrumentationVersion(),
		SchemaURL: c.SchemaURL(),
	}
	t, ok := p.namedTracer[is]
	if !ok {
		t = &logger{
			provider:             p,
			instrumentationScope: is,
		}
		p.namedTracer[is] = t
		global.Info("Logger created", "name", name, "version", c.InstrumentationVersion(), "schemaURL", c.SchemaURL())
	}
	return t
}

func WithBatcher(e LogExporter, opts ...BatchLogRecordProcessorOption) LoggerProviderOption {
	return WithLogRecordProcessor(NewBatchLogProcessorProcessor(e, opts...))
}

func WithLogRecordProcessor(lrp LogRecordProcessor) LoggerProviderOption {
	return loggerProviderOptionFunc(func(cfg loggerProviderConfig) loggerProviderConfig {
		cfg.processors = append(cfg.processors, lrp)
		return cfg
	})
}

func WithResource(r *resource.Resource) LoggerProviderOption {
	return loggerProviderOptionFunc(func(cfg loggerProviderConfig) loggerProviderConfig {
		var err error
		cfg.resource, err = resource.Merge(resource.Environment(), r)
		if err != nil {
			otel.Handle(err)
		}
		return cfg
	})
}
