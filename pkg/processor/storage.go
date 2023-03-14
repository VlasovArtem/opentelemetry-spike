package main

import (
	"context"
	"errors"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type data struct {
	name string
}

var inmemory = make(map[string]data)

func insertData(parentCtx context.Context, name string) error {
	log.Info().Msg("Inserting data")

	if _, ok := inmemory[name]; ok {
		err := errors.New("data already exists")
		return err
	}

	span := trace.SpanFromContext(parentCtx)
	span.AddEvent(
		"Inserting data",
		trace.WithAttributes(
			attribute.String("name", name),
		),
	)

	inmemory[name] = data{name: name}
	return nil
}
