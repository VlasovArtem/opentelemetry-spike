package opentelemetry

import (
	"github.com/rs/zerolog/log"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
)

func EnableGlobalLogging() func() error {
	log.Info().Msg("Enabling global logging")

	logger := otelzap.New(zap.NewExample(), otelzap.WithMinLevel(zap.DebugLevel), otelzap.WithTraceIDField(true))

	undo := otelzap.ReplaceGlobals(logger)

	otelzap.L().Info("replaced zap's global loggers")

	return func() error {
		undo()
		return logger.Sync()
	}
}
