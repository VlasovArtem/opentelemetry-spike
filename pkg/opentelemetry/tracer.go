package opentelemetry

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.18.0"
	"log"
	"os"
)

func EnableGlobalFileTracer(serviceName string) func(ctx context.Context) error {
	f, err := os.Create(serviceName + ".txt")
	if err != nil {
		log.Fatal(err)
	}

	exporter, err := NewWriterExporter(f)

	if err != nil {
		log.Fatal(err)
	}

	cleanup := initTracer(serviceName, exporter)

	return func(ctx context.Context) error {
		f.Close()
		return cleanup(ctx)
	}
}

func EnabledGlobalGrpcTracer(serviceName string, url string, insecure bool) func(ctx context.Context) error {
	exporter, err := NewGrpcExporter(url, insecure)

	if err != nil {
		log.Fatal(err)
	}

	cleanup := initTracer(serviceName, exporter)

	return cleanup
}

func initTracer(serviceName string, exporter sdktrace.SpanExporter) func(ctx context.Context) error {
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(newResource(serviceName)),
	)

	otel.SetTracerProvider(tp)

	return tp.Shutdown
}

func newResource(serviceName string) *resource.Resource {
	r, _ := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(serviceName),
			semconv.ServiceVersion("v0.1.0"),
			attribute.String("environment", "demo"),
		),
	)
	return r
}
