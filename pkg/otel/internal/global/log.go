package global

import (
	"spike-go-opentelemetry-logging/pkg/otel/log"
	"sync"
	"sync/atomic"
)

type il struct {
	name    string
	version string
}

type loggerProvider struct {
	mtx      sync.Mutex
	loggers  map[il]*logger
	delegate log.LoggerProvider
}

type logger struct {
	name     string
	opts     []log.LoggerOption
	provider *loggerProvider

	delegate atomic.Value
}

func (l *loggerProvider) Logger(name string, opts ...log.LoggerOption) log.Logger {
	l.mtx.Lock()
	defer l.mtx.Unlock()

	if l.delegate != nil {
		return l.delegate.Logger(name, opts...)
	}

	// At this moment it is guaranteed that no sdk is installed, save the tracer in the loggers map.

	c := log.NewLoggerConfig(opts...)
	key := il{
		name:    name,
		version: c.InstrumentationVersion(),
	}

	if l.loggers == nil {
		l.loggers = make(map[il]*logger)
	}

	if val, ok := l.loggers[key]; ok {
		return val
	}

	t := &logger{name: name, opts: opts, provider: l}
	l.loggers[key] = t
	return t
}

func (l *logger) Emit(record log.ReadableLogRecord) {
	delegate := l.delegate.Load()
	if delegate != nil {
		delegate.(log.Logger).Emit(record)
	}
}

func (l *logger) NewReadWriteLogRecord() log.ReadWriteLogRecord {
	delegate := l.delegate.Load()
	if delegate != nil {
		return delegate.(log.Logger).NewReadWriteLogRecord()
	}
	return nil
}

func (l *loggerProvider) setDelegate(provider log.LoggerProvider) {
	l.mtx.Lock()
	defer l.mtx.Unlock()

	l.delegate = provider

	if len(l.loggers) == 0 {
		return
	}

	for _, t := range l.loggers {
		t.setDelegate(provider)
	}

	l.loggers = nil
}

func (l *logger) setDelegate(provider log.LoggerProvider) {
	l.delegate.Store(provider.Logger(l.name, l.opts...))
}
