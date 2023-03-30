package log

import "go.opentelemetry.io/otel/attribute"

// LoggerConfig is a group of options for a Logger.
type LoggerConfig struct {
	instrumentationVersion string
	// Schema URL of the telemetry emitted by the Logger.
	schemaURL string
	attrs     attribute.Set
}

// InstrumentationVersion returns the version of the library providing instrumentation.
func (t *LoggerConfig) InstrumentationVersion() string {
	return t.instrumentationVersion
}

// InstrumentationAttributes returns the attributes associated with the library
// providing instrumentation.
func (t *LoggerConfig) InstrumentationAttributes() attribute.Set {
	return t.attrs
}

// SchemaURL returns the Schema URL of the telemetry emitted by the Tracer.
func (t *LoggerConfig) SchemaURL() string {
	return t.schemaURL
}

func NewLoggerConfig(options ...LoggerOption) LoggerConfig {
	var config LoggerConfig
	for _, option := range options {
		config = option.apply(config)
	}
	return config
}

type LoggerOption interface {
	apply(config LoggerConfig) LoggerConfig
}
