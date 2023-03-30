package log

import (
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/trace"
	"time"
)

type TraceContextConfig struct {
	TraceID trace.TraceID
	SpanID  trace.SpanID
}

type TraceContext struct {
	traceID trace.TraceID
	spanID  trace.SpanID
}

func NewTraceContext(config TraceContextConfig) TraceContext {
	return TraceContext{
		traceID: config.TraceID,
		spanID:  config.SpanID,
	}
}

func (t TraceContext) IsValid() bool {
	return t.traceID.IsValid() && t.spanID.IsValid()
}

func (t TraceContext) WithTraceID(id trace.TraceID) TraceContext {
	return TraceContext{
		traceID: id,
		spanID:  t.spanID,
	}
}

func (t TraceContext) WithSpanID(id trace.SpanID) TraceContext {
	return TraceContext{
		traceID: t.traceID,
		spanID:  id,
	}
}

func (t TraceContext) TraceID() trace.TraceID {
	return t.traceID
}

func (t TraceContext) SpanID() trace.SpanID {
	return t.spanID
}

type Logger interface {
	Emit(record ReadableLogRecord)
	NewReadWriteLogRecord() ReadWriteLogRecord
}

type LoggerProvider interface {
	Logger(name string, options ...LoggerOption) Logger
}

type SeverityNumber int32

const (
	// UNSPECIFIED is the default SeverityNumber, it MUST NOT be used.
	SeverityNumber_SEVERITY_NUMBER_UNSPECIFIED SeverityNumber = 0
	SeverityNumber_SEVERITY_NUMBER_TRACE       SeverityNumber = 1
	SeverityNumber_SEVERITY_NUMBER_TRACE2      SeverityNumber = 2
	SeverityNumber_SEVERITY_NUMBER_TRACE3      SeverityNumber = 3
	SeverityNumber_SEVERITY_NUMBER_TRACE4      SeverityNumber = 4
	SeverityNumber_SEVERITY_NUMBER_DEBUG       SeverityNumber = 5
	SeverityNumber_SEVERITY_NUMBER_DEBUG2      SeverityNumber = 6
	SeverityNumber_SEVERITY_NUMBER_DEBUG3      SeverityNumber = 7
	SeverityNumber_SEVERITY_NUMBER_DEBUG4      SeverityNumber = 8
	SeverityNumber_SEVERITY_NUMBER_INFO        SeverityNumber = 9
	SeverityNumber_SEVERITY_NUMBER_INFO2       SeverityNumber = 10
	SeverityNumber_SEVERITY_NUMBER_INFO3       SeverityNumber = 11
	SeverityNumber_SEVERITY_NUMBER_INFO4       SeverityNumber = 12
	SeverityNumber_SEVERITY_NUMBER_WARN        SeverityNumber = 13
	SeverityNumber_SEVERITY_NUMBER_WARN2       SeverityNumber = 14
	SeverityNumber_SEVERITY_NUMBER_WARN3       SeverityNumber = 15
	SeverityNumber_SEVERITY_NUMBER_WARN4       SeverityNumber = 16
	SeverityNumber_SEVERITY_NUMBER_ERROR       SeverityNumber = 17
	SeverityNumber_SEVERITY_NUMBER_ERROR2      SeverityNumber = 18
	SeverityNumber_SEVERITY_NUMBER_ERROR3      SeverityNumber = 19
	SeverityNumber_SEVERITY_NUMBER_ERROR4      SeverityNumber = 20
	SeverityNumber_SEVERITY_NUMBER_FATAL       SeverityNumber = 21
	SeverityNumber_SEVERITY_NUMBER_FATAL2      SeverityNumber = 22
	SeverityNumber_SEVERITY_NUMBER_FATAL3      SeverityNumber = 23
	SeverityNumber_SEVERITY_NUMBER_FATAL4      SeverityNumber = 24
)

func (s SeverityNumber) String() string {
	switch s {
	case SeverityNumber_SEVERITY_NUMBER_UNSPECIFIED:
		return "UNSPECIFIED"
	case SeverityNumber_SEVERITY_NUMBER_TRACE:
		return "TRACE"
	case SeverityNumber_SEVERITY_NUMBER_TRACE2:
		return "TRACE2"
	case SeverityNumber_SEVERITY_NUMBER_TRACE3:
		return "TRACE3"
	case SeverityNumber_SEVERITY_NUMBER_TRACE4:
		return "TRACE4"
	case SeverityNumber_SEVERITY_NUMBER_DEBUG:
		return "DEBUG"
	case SeverityNumber_SEVERITY_NUMBER_DEBUG2:
		return "DEBUG2"
	case SeverityNumber_SEVERITY_NUMBER_DEBUG3:
		return "DEBUG3"
	case SeverityNumber_SEVERITY_NUMBER_DEBUG4:
		return "DEBUG4"
	case SeverityNumber_SEVERITY_NUMBER_INFO:
		return "INFO"
	case SeverityNumber_SEVERITY_NUMBER_INFO2:
		return "INFO2"
	case SeverityNumber_SEVERITY_NUMBER_INFO3:
		return "INFO3"
	case SeverityNumber_SEVERITY_NUMBER_INFO4:
		return "INFO4"
	case SeverityNumber_SEVERITY_NUMBER_WARN:
		return "WARN"
	case SeverityNumber_SEVERITY_NUMBER_WARN2:
		return "WARN2"
	case SeverityNumber_SEVERITY_NUMBER_WARN3:
		return "WARN3"
	case SeverityNumber_SEVERITY_NUMBER_WARN4:
		return "WARN4"
	case SeverityNumber_SEVERITY_NUMBER_ERROR:
		return "ERROR"
	case SeverityNumber_SEVERITY_NUMBER_ERROR2:
		return "ERROR2"
	case SeverityNumber_SEVERITY_NUMBER_ERROR3:
		return "ERROR3"
	case SeverityNumber_SEVERITY_NUMBER_ERROR4:
		return "ERROR4"
	case SeverityNumber_SEVERITY_NUMBER_FATAL:
		return "FATAL"
	case SeverityNumber_SEVERITY_NUMBER_FATAL2:
		return "FATAL2"
	case SeverityNumber_SEVERITY_NUMBER_FATAL3:
		return "FATAL3"
	case SeverityNumber_SEVERITY_NUMBER_FATAL4:
		return "FATAL4"
	default:
		return "UNSPECIFIED"
	}
}

type ReadableLogRecord interface {
	Timestamp() time.Time
	ObservedTimestamp() time.Time
	SeverityNumber() SeverityNumber
	SeverityText() string
	Body() attribute.Value
	Attributes() []attribute.KeyValue
	Context() TraceContext
	Resource() *resource.Resource
	InstrumentationScope() instrumentation.Scope
}

type ReadWriteLogRecord interface {
	ReadableLogRecord
	SetTimestamp(time.Time)
	SetObservedTimestamp(time.Time)
	SetSeverityNumber(SeverityNumber)
	SetSeverityText(string)
	SetBody(value attribute.Value)
	SetAttributes(attributes ...attribute.KeyValue)
	SetContext(TraceContext)
}
