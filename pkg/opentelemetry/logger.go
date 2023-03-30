package opentelemetry

import (
	"context"
	"github.com/rs/zerolog/log"
	"spike-go-opentelemetry-logging/pkg/otel"
	sdklog "spike-go-opentelemetry-logging/pkg/otel/sdk/log"
)

func EnabledGlobalGrpcLogger(serviceName string, url string, insecure bool) func(ctx context.Context) error {
	exporter, err := NewGrpcLogRecordsExporter(url, insecure)

	if err != nil {
		log.Fatal().Err(err).Msg("failed to create exporter")
	}

	cleanup := initLogger(serviceName, exporter)

	return cleanup
}

func initLogger(serviceName string, exporter sdklog.LogExporter) func(ctx context.Context) error {
	lp := sdklog.NewLoggerProvider(
		sdklog.WithBatcher(exporter),
		sdklog.WithResource(newResource(serviceName)),
	)

	otel.SetLoggerProvider(lp)

	return lp.Shutdown
}
