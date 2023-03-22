package main

import (
	"context"
	"encoding/json"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"spike-go-opentelemetry-logging/pkg/common"
	"spike-go-opentelemetry-logging/pkg/opentelemetry"
)

func receiveRequest(conn *kafka.Conn) (request insertDataRequest, ctx context.Context, err error) {
	// We need to decide should we create a new span for each message or not
	message, err := conn.ReadMessage(1e3)
	if err != nil {
		log.Error().Err(err)
		return
	}

	ctx = otel.GetTextMapPropagator().Extract(context.Background(), opentelemetry.NewMessageCarrier(&message))

	ctx, span := tracer.Start(ctx, "kafka.consumer",
		trace.WithSpanKind(trace.SpanKindConsumer),
		common.CreateRequiredKafkaOtelConsumerAttributes(
			common.GlobalOpts.Kafka.Topic,
			0,
		),
	)

	defer span.End()

	err = json.Unmarshal(message.Value, &request)

	if err != nil {
		return
	}

	otelzap.Ctx(ctx).Info("Received message", zap.String("name", request.Name), zap.Int("random", request.Random))

	return
}
