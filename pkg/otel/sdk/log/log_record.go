package log

import (
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/resource"
	"spike-go-opentelemetry-logging/pkg/otel/log"
	"sync"
	"time"
)

type logRecord struct {
	mu                sync.Mutex
	timestamp         time.Time
	observedTimestamp time.Time
	severityNumber    log.SeverityNumber
	severityText      string
	body              attribute.Value
	attributes        []attribute.KeyValue
	context           log.TraceContext
	logger            *logger
}

func (l *logRecord) Timestamp() time.Time {
	if l == nil {
		return time.Time{}
	}
	return l.timestamp
}

func (l *logRecord) ObservedTimestamp() time.Time {
	if l == nil {
		return time.Time{}
	}
	return l.observedTimestamp
}

func (l *logRecord) SeverityNumber() log.SeverityNumber {
	if l == nil {
		return log.SeverityNumber_SEVERITY_NUMBER_UNSPECIFIED
	}
	return l.severityNumber
}

func (l *logRecord) SeverityText() string {
	if l == nil {
		return ""
	}
	if l.severityText == "" {
		return l.severityNumber.String()
	}
	return l.severityText
}

func (l *logRecord) Body() attribute.Value {
	if l == nil {
		return attribute.Value{}
	}
	return l.body
}

func (l *logRecord) Attributes() []attribute.KeyValue {
	if l == nil {
		return nil
	}
	return l.attributes
}

func (l *logRecord) Context() log.TraceContext {
	if l == nil {
		return log.TraceContext{}
	}
	return l.context
}

func (l *logRecord) Resource() *resource.Resource {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.logger.provider.resource
}

func (l *logRecord) InstrumentationScope() instrumentation.Scope {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.logger.instrumentationScope
}

type mutableLogRecord struct {
	*logRecord
}

func NewReadWriteLogRecord(logger *logger) log.ReadWriteLogRecord {
	return &mutableLogRecord{
		logRecord: &logRecord{
			timestamp: time.Now(),
			logger: logger,
		},
	}
}

func (m *mutableLogRecord) SetTimestamp(t time.Time) {
	if m == nil {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.timestamp = t
}

func (m *mutableLogRecord) SetObservedTimestamp(t time.Time) {
	if m == nil {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.observedTimestamp = t
}

func (m *mutableLogRecord) SetSeverityNumber(number log.SeverityNumber) {
	if m == nil {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.severityNumber = number
}

func (m *mutableLogRecord) SetSeverityText(severityText string) {
	if m == nil {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.severityText = severityText
}

func (m *mutableLogRecord) SetBody(body attribute.Value) {
	if m == nil {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.body = body
}

func (m *mutableLogRecord) SetAttributes(attributes ...attribute.KeyValue) {
	if m == nil {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.attributes = attributes
}

func (m *mutableLogRecord) SetContext(context log.TraceContext) {
	if m == nil {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.context = context
}
