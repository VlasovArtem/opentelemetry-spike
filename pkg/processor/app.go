package main

import (
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"spike-go-opentelemetry-logging/pkg/common"
	"sync"
	"time"
)

const serviceName = "processor-app"

type insertDataRequest struct {
	Name   string `json:"name"`
	Random int    `json:"random"`
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

	log.Info().Msg("Creating kafka connection")
	log.Info().Msgf("Kafka Configuration %v", common.GlobalOpts.Kafka)

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

	lock := sync.Mutex{}

	for {
		if lock.TryLock() {
			request, ctx, err := receiveRequest(connection)
			if err != nil {
				log.Error().Err(err).Msg("Error receiving request")
			} else {
				parent, span := tracer.Start(ctx, "dataProcessing",
					trace.WithSpanKind(trace.SpanKindServer),
					trace.WithAttributes(
						attribute.String("name", request.Name),
						attribute.Int("random", request.Random),
					),
				)
				err := insertData(parent, request)
				if err == nil {
					execute(parent, request)
				}
				span.End()
			}
			lock.Unlock()
		}
		time.Sleep(1 * time.Second)
	}
}
