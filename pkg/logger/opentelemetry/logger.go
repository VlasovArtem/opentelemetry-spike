package opentelemetry

import (
	"context"
	"fmt"
	zlog "github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"spike-go-opentelemetry-logging/pkg/logger"
	"spike-go-opentelemetry-logging/pkg/otel/log"
	"sync"
)

type otelProxyLogger struct {
	mu     sync.Mutex
	logger log.Logger
	ctx    context.Context
}

func NewOtelProxyLogger(otelLogger log.Logger) logger.Logger {
	return &otelProxyLogger{logger: otelLogger}
}

func (o *otelProxyLogger) Debug() logger.Event {
	record := o.logger.NewReadWriteLogRecord()
	record.SetSeverityNumber(log.SeverityNumber_SEVERITY_NUMBER_DEBUG)
	return &otelProxyEvent{
		logger: o,
		record: record,
	}
}

func (o *otelProxyLogger) Info() logger.Event {
	record := o.logger.NewReadWriteLogRecord()
	record.SetSeverityNumber(log.SeverityNumber_SEVERITY_NUMBER_INFO)
	return &otelProxyEvent{
		logger: o,
		record: record,
	}
}

func (o *otelProxyLogger) Warn() logger.Event {
	record := o.logger.NewReadWriteLogRecord()
	record.SetSeverityNumber(log.SeverityNumber_SEVERITY_NUMBER_WARN)
	return &otelProxyEvent{
		logger: o,
		record: record,
	}
}

func (o *otelProxyLogger) Error() logger.Event {
	record := o.logger.NewReadWriteLogRecord()
	record.SetSeverityNumber(log.SeverityNumber_SEVERITY_NUMBER_ERROR)
	return &otelProxyEvent{
		logger: o,
		record: record,
	}
}

func (o *otelProxyLogger) Fatal() logger.Event {
	record := o.logger.NewReadWriteLogRecord()
	record.SetSeverityNumber(log.SeverityNumber_SEVERITY_NUMBER_FATAL)
	return &otelProxyEvent{
		logger: o,
		record: record,
	}
}

func (o *otelProxyLogger) Contextual(ctx context.Context) logger.Logger {
	o.ctx = ctx
	return o
}

func (o *otelProxyLogger) emit(event otelProxyEvent) {
	o.mu.Lock()
	defer o.mu.Unlock()
	record := event.record
	if o.ctx == nil {
		zlog.Error().Msg("context not set")
		return
	}
	span := trace.SpanFromContext(o.ctx)
	if span == nil {
		zlog.Error().Msg("span not found in context")
		return
	}
	spanContext := span.SpanContext()
	traceContext := log.NewTraceContext(
		log.TraceContextConfig{
			TraceID: spanContext.TraceID(),
			SpanID:  spanContext.SpanID(),
		},
	)
	if !traceContext.IsValid() {
		zlog.Error().Msg("trace context not valid")
		return
	}
	record.SetContext(traceContext)
	o.logger.Emit(record)
}

type otelProxyEvent struct {
	logger *otelProxyLogger
	record log.ReadWriteLogRecord
}

func (o *otelProxyEvent) Msg(msg string) {
	o.record.SetBody(attribute.StringValue(msg))
	o.logger.emit(*o)
}

func (o *otelProxyEvent) Msgf(format string, args ...any) {
	o.record.SetBody(attribute.StringValue(fmt.Sprintf(format, args...)))
	o.logger.emit(*o)
}

func (o *otelProxyEvent) MsgE(msg string, err error) {
	o.record.SetBody(attribute.StringValue(msg))
	o.record.SetAttributes(attribute.String("error", err.Error()))
	o.logger.emit(*o)
}
