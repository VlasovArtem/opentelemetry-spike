package main

import (
	"github.com/gin-gonic/gin"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"net/http"
)

func initHandler(r *gin.Engine) {
	r.Use(otelgin.Middleware(serviceName))

	r.GET("/ping", ping())
	r.POST("/data", insertData())
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
		var request insertDataRequest
		err := context.Bind(&request)
		if err != nil {
			otelzap.Ctx(context).Error("Error binding request", zap.Error(err))
			context.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid request",
			})
		} else {
			span.AddEvent("Inserting data", trace.WithAttributes(
				attribute.String("name", request.Name),
			))
			err := sendMessage(request, ctx)
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
