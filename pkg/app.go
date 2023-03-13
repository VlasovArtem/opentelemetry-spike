package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/jessevdk/go-flags"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"log"
	"net/http"
)

type opts struct {
	ExporterType string `short:"t" long:"type" required:"true" description:"type of exporter: file or grpc" choice:"file" choice:"grpc"`
	Collector    struct {
		Url      string `short:"u" long:"url" description:"collector url"`
		Insecure bool   `short:"i" long:"insecure" description:"collector grpc insecure"`
	} `group:"collector" namespace:"collector"`
}

func main() {
	opts := opts{}
	_, err := flags.Parse(&opts)

	if err != nil {
		log.Fatal(err)
	}

	switch opts.ExporterType {
	case "file":
		defer initFileTracer()(context.Background())
		log.Println("Starting application with file exporter. Check 'traces.txt' file for traces")
	case "grpc":
		defer initGrpcTracer(opts.Collector.Url, opts.Collector.Insecure)(context.Background())
		log.Println("Starting application with grpc exporter.")
	default:
		log.Fatal("Invalid type")
	}

	defer initGlobalLogging()()

	r := gin.Default()

	r.Use(otelgin.Middleware("application"))

	r.GET("/ping", ping())
	r.POST("/data", insertData())
	r.Run(":8080")
}

func ping() gin.HandlerFunc {
	return func(c *gin.Context) {
		otelzap.Ctx(c).Info("Ping")
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	}
}

type insertDataRequest struct {
	Name string `json:"name"`
}

func insertData() gin.HandlerFunc {
	return func(context *gin.Context) {
		ctx := context.Request.Context()
		span := trace.SpanFromContext(ctx)
		span.AddEvent("Inserting data with name")
		var request insertDataRequest
		err := context.Bind(&request)
		if err != nil {
			otelzap.Ctx(context).Error("Error binding request", zap.Error(err))
			context.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid request",
			})
		} else {
			err := InsertData(request.Name)
			if err != nil {
				otelzap.Ctx(ctx).Error("Error inserting data", zap.Error(err))
				context.JSON(http.StatusInternalServerError, gin.H{
					"message": "Error inserting data",
				})
			} else {
				otelzap.Ctx(ctx).Info("Data inserted")
				context.JSON(http.StatusOK, gin.H{
					"message": "Data inserted",
				})
			}
		}
	}
}
