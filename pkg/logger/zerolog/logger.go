package zerolog

import (
	"context"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"spike-go-opentelemetry-logging/pkg/logger"
)

type zerologProxy struct{}

func NewZerologProxyLogger() logger.Logger {
	return &zerologProxy{}
}

func (z *zerologProxy) Debug() logger.Event {
	return &zerologEventProxy{event: log.Debug()}
}

func (z *zerologProxy) Info() logger.Event {
	return &zerologEventProxy{event: log.Info()}
}

func (z *zerologProxy) Warn() logger.Event {
	return &zerologEventProxy{event: log.Warn()}
}

func (z *zerologProxy) Error() logger.Event {
	return &zerologEventProxy{event: log.Error()}
}

func (z *zerologProxy) Fatal() logger.Event {
	return &zerologEventProxy{event: log.Fatal()}
}

func (z *zerologProxy) Contextual(ctx context.Context) logger.Logger {
	return &zerologProxyContextual{
		logger: log.Ctx(ctx),
	}
}

type zerologProxyContextual struct {
	logger *zerolog.Logger
}

func (z *zerologProxyContextual) Debug() logger.Event {
	return &zerologEventProxy{event: z.logger.Debug()}
}

func (z *zerologProxyContextual) Info() logger.Event {
	return &zerologEventProxy{event: z.logger.Info()}
}

func (z *zerologProxyContextual) Warn() logger.Event {
	return &zerologEventProxy{event: z.logger.Warn()}
}

func (z *zerologProxyContextual) Error() logger.Event {
	return &zerologEventProxy{event: z.logger.Error()}
}

func (z *zerologProxyContextual) Fatal() logger.Event {
	return &zerologEventProxy{event: z.logger.Fatal()}
}

func (z *zerologProxyContextual) Contextual(ctx context.Context) logger.Logger {
	return &zerologProxyContextual{
		logger: log.Ctx(ctx),
	}
}

type zerologEventProxy struct {
	event *zerolog.Event
}

func (z *zerologEventProxy) Msg(msg string) {
	z.event.Msg(msg)
}

func (z *zerologEventProxy) Msgf(format string, args ...any) {
	z.event.Msgf(format, args)
}

func (z *zerologEventProxy) MsgE(msg string, err error) {
	z.event.Err(err).Msg(msg)
}
