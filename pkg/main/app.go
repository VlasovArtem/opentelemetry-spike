package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/trace"

	"spike-go-opentelemetry-logging/pkg/common"
)

const serviceName = "handler-app"

var tracer trace.Tracer
var meter metric.Meter

func main() {
	err := common.ParseGlobalOpts()
	if err != nil {
		log.Fatal(err)
	}

	defer common.EnableTelemetry(serviceName)()

	tracer = otel.GetTracerProvider().Tracer(serviceName)
	meter = global.MeterProvider().Meter(serviceName)

	r := gin.Default()

	initHandler(r)

	r.Run(":8080")
}
