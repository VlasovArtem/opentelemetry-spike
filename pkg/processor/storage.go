package main

import (
	"context"
	"errors"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

var inmemory = make(map[string]insertDataRequest)

func insertData(parentCtx context.Context, request insertDataRequest) error {
	ctx, span := tracer.Start(parentCtx, "insertData",
		trace.WithSpanKind(trace.SpanKindServer),
		trace.WithAttributes(
			attribute.String("name", request.Name),
		),
	)
	defer span.End()
	otelzap.Ctx(ctx).ZapLogger().Info("Inserting data", zap.String("name", request.Name))

	if _, ok := inmemory[request.Name]; ok {
		err := errors.New("data already exists")
		otelzap.L().Error("Data already exists", zap.String("name", request.Name))
		span.RecordError(err)
		return err
	}

	inmemory[request.Name] = request
	return nil
}
