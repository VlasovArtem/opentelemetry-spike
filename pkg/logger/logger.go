package logger

import (
	"context"
)

type Event interface {
	Msg(msg string)
	Msgf(format string, args ...any)
	MsgE(msg string, err error)
}

type Logger interface {
	Debug() Event
	Info() Event
	Warn() Event
	Error() Event
	Fatal() Event
	Contextual(ctx context.Context) Logger
}
