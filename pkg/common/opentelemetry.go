package common

import (
	"context"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
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

func otelKafkaDynamicAttributes(topic string, partition int, messageId string) []attribute.KeyValue {
	defaultAttributes := otelKafkaDefaultAttributes()

	return append(defaultAttributes,
		semconv.MessagingDestinationKey.String(topic),
		semconv.MessagingMessageIDKey.String(messageId),
		semconv.MessagingKafkaPartitionKey.Int(partition),
	)
}

func otelKafkaDefaultAttributes() []attribute.KeyValue {
	return []attribute.KeyValue{
		semconv.MessagingSystemKey.String("kafka"),
		semconv.MessagingDestinationKindTopic,
	}
}

func CreateRequiredKafkaOtelConsumerAttributes(topic string, partition int, messageId string) trace.SpanStartOption {
	attributes := otelKafkaDynamicAttributes(topic, partition, messageId)
	attributes = append(attributes, semconv.MessagingOperationReceive)

	return trace.WithAttributes(attributes...)
}

func CreateRequiredKafkaOtelProducerAttributes(topic string, partition int, messageId string) trace.SpanStartOption {
	return trace.WithAttributes(otelKafkaDynamicAttributes(topic, partition, messageId)...)
}
