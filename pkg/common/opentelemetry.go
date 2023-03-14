package common

import (
	"context"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"spike-go-opentelemetry-logging/pkg/opentelemetry"
)

func EnableTelemetry(serviceName string) func() error {
	log.Info().Msg("Enabling telemetry")

	var tracerCloseFunc func(ctx context.Context) error
	switch GlobalOpts.ExporterType {
	case "file":
		tracerCloseFunc = opentelemetry.EnableGlobalFileTracer(serviceName)
		log.Info().Msg("Starting application with file exporter. Check 'traces.txt' file for traces")
	case "grpc":
		tracerCloseFunc = opentelemetry.EnabledGlobalGrpcTracer(
			serviceName,
			GlobalOpts.Collector.Url,
			GlobalOpts.Collector.Insecure,
		)
		log.Info().Msg("Starting application with grpc exporter.")
	default:
		log.Fatal().Msg("Invalid type")
	}

	loggingCloseFunc := opentelemetry.EnableGlobalLogging()

	return func() error {
		loggingErr := loggingCloseFunc()
		if loggingErr != nil {
			log.Error().Err(loggingErr)
		}
		tracerErr := tracerCloseFunc(context.Background())
		if tracerErr != nil {
			log.Error().Err(tracerErr)
			if loggingErr != nil {
				return errors.Wrap(loggingErr, tracerErr.Error())
			}
			return tracerErr
		} else {
			return loggingErr
		}
	}
}
