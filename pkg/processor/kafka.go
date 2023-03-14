package main

import (
	"context"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func receiveRequest(conn *kafka.Conn) (insertDataRequest, context.Context, error) {
	log.Info().Msg("Waiting for message...")

	// We need to decide should we create a new span for each message or not
	message, err := conn.ReadMessage(1e3)
	if err != nil {
		log.Error().Err(err)
		return insertDataRequest{}, nil, err
	}

	ctx, span := tracer.Start(context.TODO(), "add.data receive",
		trace.WithSpanKind(trace.SpanKindConsumer),
		trace.WithAttributes(
			attribute.String("messaging.system", "kafka"),
			attribute.String("messaging.operation", "receive"),
			attribute.String("messaging.message.conversation_id", findConversationId(message.Headers)),
		),
	)

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
