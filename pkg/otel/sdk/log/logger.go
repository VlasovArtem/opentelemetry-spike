package log

import (
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"spike-go-opentelemetry-logging/pkg/otel/log"
)

type logger struct {
	provider             *LoggerProvider
	instrumentationScope instrumentation.Scope
}

func (l *logger) Emit(r log.ReadableLogRecord) {
	sps := l.provider.logsProcessors.Load().(logRecordProcessorStates)
	for _, sp := range sps {
		sp.sp.OnEmit(r)
	}
}

func (l *logger) NewReadWriteLogRecord() log.ReadWriteLogRecord {
	return NewReadWriteLogRecord(l)
}
