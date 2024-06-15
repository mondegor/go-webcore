package mrlog

import (
	"context"
)

//go:generate mockgen -source=logger.go -destination=./mock/logger.go

type (
	// Logger - интерфейс логирования ошибок и сообщений через формирования события.
	Logger interface { //nolint:interfacebloat
		Level() Level
		WithContext(ctx context.Context) context.Context

		With() LoggerContext

		Debug() LoggerEvent
		Info() LoggerEvent
		Warn() LoggerEvent
		Error() LoggerEvent
		Fatal() LoggerEvent
		Panic() LoggerEvent
		Trace() LoggerEvent

		Printf(format string, args ...any) // поддержка стандартного интерфейса логгирования сообщений
	}

	// LoggerContext - контекст логгера для обогащения его дополнительными атрибутами.
	LoggerContext interface {
		Logger() Logger

		Str(key, value string) LoggerContext
		Bytes(key string, value []byte) LoggerContext
		Int(key string, value int) LoggerContext
		Any(key string, value any) LoggerContext
	}

	// LoggerEvent - инерфейс события, с возможностью его обогащения и отправки.
	LoggerEvent interface {
		Err(err error) LoggerEvent
		Str(key, value string) LoggerEvent
		Bytes(key string, value []byte) LoggerEvent
		Int(key string, value int) LoggerEvent
		Any(key string, value any) LoggerEvent

		Msg(message string)
		Msgf(format string, args ...any)
		MsgFunc(createMsg func() string)
		Send()
	}
)
