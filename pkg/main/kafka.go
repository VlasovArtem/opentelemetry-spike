package main

import (
	"context"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"spike-go-opentelemetry-logging/pkg/common"
)

func sendMessage(request insertDataRequest, parentCtx context.Context) error {
	conversationId := uuid.New().String()
	_, span := tracer.Start(parentCtx, "add.data publish",
		trace.WithAttributes(
			attribute.String("name", request.Name),
			attribute.String("messaging.system", "kafka"),
			attribute.String("messaging.operation", "publish"),
			attribute.String("messaging.message.conversation_id", conversationId),
		),
		trace.WithSpanKind(trace.SpanKindProducer),
	)
	defer span.End()

	newConnection := common.NewKafkaConnection(
		"tcp",
		common.GlobalOpts.Kafka.Address,
		common.GlobalOpts.Kafka.Topic,
		common.GlobalOpts.Kafka.Partition,
	)

	common.WriteMessages(newConnection, []kafka.Message{
		createMessage(request, conversationId),
	})

	err := newConnection.Close()
	if err != nil {
		otelzap.Ctx(parentCtx).Error("Error sending message to kafka", zap.Error(err))
		span.RecordError(err)
		return err
	}
	return nil
}

func createMessage(request insertDataRequest, conversationId string) kafka.Message {
	return kafka.Message{
		Key:   []byte("insert-data"),
		Value: []byte(request.Name),
		Headers: []kafka.Header{
			{
				Key:   "conversationId",
				Value: []byte(conversationId),
			},
		},
	}
}
