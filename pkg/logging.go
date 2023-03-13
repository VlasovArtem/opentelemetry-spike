package main

import (
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
)

func initGlobalLogging() func() error {
	logger := otelzap.New(zap.NewExample(), otelzap.WithMinLevel(zap.DebugLevel), otelzap.WithTraceIDField(true))

	undo := otelzap.ReplaceGlobals(logger)

	otelzap.L().Info("replaced zap's global loggers")

	return func() error {
		undo()
		return logger.Sync()
	}
}
