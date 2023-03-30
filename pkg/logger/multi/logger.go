package multi

import (
	"context"
	"github.com/rs/zerolog/log"
	"spike-go-opentelemetry-logging/pkg/logger"
	"spike-go-opentelemetry-logging/pkg/logger/zerolog"
)

type multiLogger struct {
	loggers []logger.Logger
}

func NewMultiLogger(loggers ...logger.Logger) logger.Logger {
	if len(loggers) == 0 {
		log.Warn().Msg("No loggers provided, defaulting to zerolog proxy logger")
		proxyLogger := zerolog.NewZerologProxyLogger()
		loggers = []logger.Logger{proxyLogger}
	}
	return &multiLogger{loggers: loggers}
}

func (m *multiLogger) Debug() logger.Event {
	event := multiEvent{}
	for _, l := range m.loggers {
		event.events = append(event.events, l.Debug())
	}
	return &event
}

func (m *multiLogger) Info() logger.Event {
	event := multiEvent{}
	for _, l := range m.loggers {
		event.events = append(event.events, l.Info())
	}
	return &event
}

func (m *multiLogger) Warn() logger.Event {
	event := multiEvent{}
	for _, l := range m.loggers {
		event.events = append(event.events, l.Warn())
	}
	return &event
}

func (m *multiLogger) Error() logger.Event {
	event := multiEvent{}
	for _, l := range m.loggers {
		event.events = append(event.events, l.Error())
	}
	return &event
}

func (m *multiLogger) Fatal() logger.Event {
	event := multiEvent{}
	for _, l := range m.loggers {
		event.events = append(event.events, l.Fatal())
	}
	return &event
}

func (m *multiLogger) Contextual(ctx context.Context) logger.Logger {
	loggers := make([]logger.Logger, len(m.loggers))
	for _, l := range m.loggers {
		loggers = append(loggers, l.Contextual(ctx))
	}
	return &multiLogger{loggers: loggers}
}

type multiEvent struct {
	events []logger.Event
}

func (m *multiEvent) Msg(msg string) {
	for _, e := range m.events {
		e.Msg(msg)
	}
}

func (m *multiEvent) Msgf(format string, args ...any) {
	for _, e := range m.events {
		e.Msgf(format, args...)
	}
}

func (m *multiEvent) MsgE(msg string, err error) {
	for _, e := range m.events {
		e.MsgE(msg, err)
	}
}
