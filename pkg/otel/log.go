package otel

import (
	"spike-go-opentelemetry-logging/pkg/otel/internal/global"
	"spike-go-opentelemetry-logging/pkg/otel/log"
)

func GetLoggerProvider() log.LoggerProvider {
	return global.LoggerProvider()
}

// SetLoggerProvider registers `tp` as the global trace provider.
func SetLoggerProvider(lp log.LoggerProvider) {
	global.SetLoggerProvider(lp)
}
