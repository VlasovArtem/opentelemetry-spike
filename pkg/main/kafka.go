package main

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"spike-go-opentelemetry-logging/pkg/common"
	"spike-go-opentelemetry-logging/pkg/opentelemetry"
)

func sendMessage(request insertDataRequest, parentCtx context.Context) error {
	currentContext, span := tracer.Start(parentCtx, "kafka.producer",
		common.CreateRequiredKafkaOtelProducerAttributes(
			common.GlobalOpts.Kafka.Topic,
			0,
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

	message, err := createMessage(request)
	if err != nil {
		otelzap.Ctx(parentCtx).Error("Error creating message", zap.Error(err))
		span.RecordError(err)
		return err
	}
	messageCarrier := opentelemetry.NewMessageCarrier(&message)
	otel.GetTextMapPropagator().Inject(currentContext, messageCarrier)
	common.WriteMessages(newConnection, []kafka.Message{
		message,
	})

	err = newConnection.Close()
	if err != nil {
		otelzap.Ctx(parentCtx).Error("Error sending message to kafka", zap.Error(err))
		span.RecordError(err)
		return err
	}
	return nil
}

func createMessage(request insertDataRequest) (kafka.Message, error) {
	content, err := json.Marshal(request)
	if err != nil {
		return kafka.Message{}, err
	}
	return kafka.Message{
		Key:   []byte("insert-data"),
		Value: content,
	}, nil
}
