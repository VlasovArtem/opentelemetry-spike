package main

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	sdk_trace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.18.0"
	"log"
	"os"
)

func initFileTracer() func(ctx context.Context) error {
	f, err := os.Create("traces.txt")
	if err != nil {
		log.Fatal(err)
	}

	exporter, err := newWriterExporter(f)

	if err != nil {
		log.Fatal(err)
	}

	cleanup := initTracer(exporter)

	return func(ctx context.Context) error {
		f.Close()
		return cleanup(ctx)
	}
}

func initGrpcTracer(url string, insecure bool) func(ctx context.Context) error {
	exporter, err := newGrpcExporter(url, insecure)

	if err != nil {
		log.Fatal(err)
	}

	cleanup := initTracer(exporter)

	return cleanup
}

func initTracer(exporter trace.SpanExporter) func(ctx context.Context) error {
	tp := sdk_trace.NewTracerProvider(
		sdk_trace.WithBatcher(exporter),
		sdk_trace.WithResource(newResource()),
	)

	otel.SetTracerProvider(tp)

	return tp.Shutdown
}

func newResource() *resource.Resource {
	r, _ := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("application"),
			semconv.ServiceVersion("v0.1.0"),
			attribute.String("environment", "demo"),
		),
	)
	return r
}
