package main

import (
	"context"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"spike-go-opentelemetry-logging/pkg/common"
	"spike-go-opentelemetry-logging/pkg/opentelemetry"
)

func receiveRequest(conn *kafka.Conn) (insertDataRequest, context.Context, error) {
	// We need to decide should we create a new span for each message or not
	message, err := conn.ReadMessage(1e3)
	if err != nil {
		log.Error().Err(err)
		return insertDataRequest{}, nil, err
	}

	ctx := otel.GetTextMapPropagator().Extract(context.Background(), opentelemetry.NewMessageCarrier(&message))

	_, span := tracer.Start(ctx, "kafka.consumer",
		trace.WithSpanKind(trace.SpanKindConsumer),
		common.CreateRequiredKafkaOtelConsumerAttributes(
			common.GlobalOpts.Kafka.Topic,
			0,
			findConversationId(message.Headers),
		),
	)

	defer span.End()

	request := insertDataRequest{
		name: string(message.Value),
	}

	log.Info().Msgf("Received message: %s", request.name)

	span.AddEvent("received message", trace.WithAttributes(attribute.String("name", request.name)))

	return insertDataRequest{}, ctx, err
}

func findConversationId(headers []kafka.Header) string {
	for _, header := range headers {
		if header.Key == "conversationId" {
			return string(header.Value)
		}
	}
	return ""
}
