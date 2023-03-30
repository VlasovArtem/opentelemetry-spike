package log

import (
	"context"
	"spike-go-opentelemetry-logging/pkg/otel/log"
	"sync"
)

type LogRecordProcessor interface {
	OnEmit(l log.ReadableLogRecord)
	Shutdown(ctx context.Context) error
	ForceFlush(ctx context.Context) error
}

type logRecordProcessorState struct {
	sp    LogRecordProcessor
	state *sync.Once
}

func newLogRecordProcessorState(sp LogRecordProcessor) *logRecordProcessorState {
	return &logRecordProcessorState{sp: sp, state: &sync.Once{}}
}

type logRecordProcessorStates []*logRecordProcessorState
