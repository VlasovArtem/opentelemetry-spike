package global

import (
	"github.com/pkg/errors"
	"spike-go-opentelemetry-logging/pkg/otel/log"
	"sync"
	"sync/atomic"
)

type tracerProviderHolder struct {
	lp log.LoggerProvider
}

var (
	globalLogger      = defaultLoggerValue()
	delegateTraceOnce sync.Once
)

func LoggerProvider() log.LoggerProvider {
	return globalLogger.Load().(tracerProviderHolder).lp
}

func SetLoggerProvider(lp log.LoggerProvider) {
	current := LoggerProvider()

	if _, cOk := current.(*loggerProvider); cOk {
		if _, tpOk := lp.(*loggerProvider); tpOk && current == lp {
			// Do not assign the default delegating TracerProvider to delegate
			// to itself.
			Error(
				errors.New("no delegate configured in tracer provider"),
				"Setting tracer provider to it's current value. No delegate will be configured",
			)
			return
		}
	}

	delegateTraceOnce.Do(func() {
		if def, ok := current.(*loggerProvider); ok {
			def.setDelegate(lp)
		}
	})
	globalLogger.Store(tracerProviderHolder{lp: lp})
}

func defaultLoggerValue() *atomic.Value {
	v := &atomic.Value{}
	v.Store(tracerProviderHolder{lp: &loggerProvider{}})
	return v
}
