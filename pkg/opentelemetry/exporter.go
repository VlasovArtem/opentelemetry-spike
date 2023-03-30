package opentelemetry

import (
	"context"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc/credentials"
	"io"
	"spike-go-opentelemetry-logging/pkg/otel/exporters/otlp/otlplog"
	"spike-go-opentelemetry-logging/pkg/otel/exporters/otlp/otlplog/otlploggrpc"
)

func NewWriterExporter(w io.Writer) (trace.SpanExporter, error) {
	return stdouttrace.New(
		stdouttrace.WithWriter(w),
		// Use human-readable output.
		stdouttrace.WithPrettyPrint(),
		// Do not print timestamps for the demo.
		stdouttrace.WithoutTimestamps(),
	)
}

func NewGrpcTraceExporter(collectorURL string, insecure bool) (trace.SpanExporter, error) {
	secureOption := otlptracegrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
	if insecure {
		secureOption = otlptracegrpc.WithInsecure()
	}

	return otlptrace.New(
		context.Background(),
		otlptracegrpc.NewClient(
			secureOption,
			otlptracegrpc.WithEndpoint(collectorURL),
		),
	)
}

func NewGrpcLogRecordsExporter(collectorURL string, insecure bool) (otlplog.LogExporter, error) {
	secureOption := otlploggrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
	if insecure {
		secureOption = otlploggrpc.WithInsecure()
	}

	return otlplog.New(
		context.Background(),
		otlploggrpc.NewClient(
			secureOption,
			otlploggrpc.WithEndpoint(collectorURL),
		),
	)
}
