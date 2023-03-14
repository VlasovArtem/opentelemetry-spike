package main

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"log"
	"spike-go-opentelemetry-logging/pkg/common"
)

const serviceName = "handler-app"

var tracer trace.Tracer

func main() {
	err := common.ParseGlobalOpts()
	if err != nil {
		log.Fatal(err)
	}

	defer common.EnableTelemetry(serviceName)()

	tracer = otel.GetTracerProvider().Tracer(serviceName)

	r := gin.Default()

	initHandler(r)

	r.Run(":8080")
}
