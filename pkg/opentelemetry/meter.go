package opentelemetry

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/sdk/metric"
)

func EnableGlobalGrpcMetric(serviceName string, url string, insecure bool) func(ctx context.Context) error {
	options := []otlpmetricgrpc.Option{
		otlpmetricgrpc.WithEndpoint(url),
	}
	if insecure {
		options = append(options, otlpmetricgrpc.WithInsecure())
	}

	exporter, err := otlpmetricgrpc.New(context.Background(), options...)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create exporter")
	}

	reader := metric.NewPeriodicReader(
		exporter,
		metric.WithInterval(15*time.Second),
	)

	provider := metric.NewMeterProvider(metric.WithReader(reader), metric.WithResource(newResource(serviceName)))

	global.SetMeterProvider(provider)

	return func(ctx context.Context) error {
		return provider.Shutdown(ctx)
	}
}
