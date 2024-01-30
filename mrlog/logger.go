package mrlog

import (
	"context"
	"log"
)

type (
	Logger interface {
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

		Printf(format string, args ...any)
	}

	LoggerContext interface {
		Logger() Logger

		CallerWithSkipFrame(count int) LoggerContext
		Str(key, value string) LoggerContext
		Bytes(key string, value []byte) LoggerContext
		Int(key string, value int) LoggerContext
		Any(key string, value any) LoggerContext
	}

	LoggerEvent interface {
		Caller(skip ...int) LoggerEvent
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

	Options struct {
		Level           Level
		JsonFormat      bool
		TimestampFormat string
		ConsoleColor    bool
		BehaviorFunc    func(err error) bool
	}

	ctxKey struct{}
)

var (
	def Logger = &DefaultLogger{
		level:  DebugLevel,
		native: log.Default(),
	}
)

func Default() Logger {
	return def
}

// SetDefault - WARNING: use only when starting the main process
func SetDefault(logger Logger) {
	if logger != nil {
		def = logger.With().Str("logger", "DEFAULT").Logger()
	}
}

func WithContext(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, ctxKey{}, logger)
}

func Ctx(ctx context.Context) Logger {
	if value, ok := ctx.Value(ctxKey{}).(Logger); ok {
		return value
	}

	return def
}
