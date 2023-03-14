package main

import (
	"github.com/rs/zerolog/log"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"spike-go-opentelemetry-logging/pkg/common"
	"time"
)

const serviceName = "processor-app"

type insertDataRequest struct {
	name string
}

var tracer trace.Tracer

func main() {
	log.Info().Msg("Starting processor app")
	log.Info().Msg("Parsing global options")
	err := common.ParseGlobalOpts()
	if err != nil {
		log.Fatal().Err(err)
	}

	defer common.EnableTelemetry(serviceName)()

	tracer = otel.GetTracerProvider().Tracer(serviceName)

	log.Info().Msg("Starting processor app loop")

	connection := common.NewKafkaConnection(
		"tcp",
		common.GlobalOpts.Kafka.Address,
		common.GlobalOpts.Kafka.Topic,
		common.GlobalOpts.Kafka.Partition,
	)

	defer func() {
		err = connection.Close()
		if err != nil {
			log.Fatal().Err(err)
		}
	}()

	for {
		request, ctx, err := receiveRequest(connection)
		if err != nil {
			log.Error().Err(err)
		} else {
			span := trace.SpanFromContext(ctx)
			err := insertData(ctx, request.name)
			if err != nil {
				otelzap.Ctx(ctx).Error("Error binding request", zap.Error(err))
				span.RecordError(err)
			}
			span.End()
		}
		time.Sleep(1 * time.Second)
	}
}
