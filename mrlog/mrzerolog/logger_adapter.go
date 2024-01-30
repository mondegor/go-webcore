package mrzerolog

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/mondegor/go-webcore/mrlog"
	"github.com/rs/zerolog"
)

// go get -u github.com/rs/zerolog

type (
	LoggerAdapter struct {
		zl    zerolog.Logger
		level mrlog.Level
	}
)

// SetDateTimeFormat - WARNING: use only when starting the main process
func SetDateTimeFormat(value string) {
	if value == "" {
		return
	}

	zerolog.TimeFieldFormat = value
}

func New(opts mrlog.Options) *LoggerAdapter {
	var out io.Writer = os.Stdout

	if !opts.JsonFormat {
		if opts.TimestampFormat == "" {
			opts.TimestampFormat = time.Kitchen
		}

		out = zerolog.ConsoleWriter{
			Out:        out,
			TimeFormat: opts.TimestampFormat,
			NoColor:    !opts.ConsoleColor,
		}
	}

	logger := zerolog.New(out)

	if opts.TimestampFormat != "" {
		logger = logger.With().Timestamp().Logger()
	}

	logger = logger.Level(zerolog.Level(opts.Level))

	return &LoggerAdapter{
		zl:    logger,
		level: opts.Level,
	}
}

func NewConsoleLogger() *LoggerAdapter {
	return New(
		mrlog.Options{
			Level:           mrlog.DebugLevel,
			JsonFormat:      false,
			TimestampFormat: time.TimeOnly,
			ConsoleColor:    true,
		},
	)
}

func NewZeroLogger(zl zerolog.Logger) *LoggerAdapter {
	level := zl.GetLevel()
	fatalLevel := zerolog.Level(mrlog.FatalLevel)

	if level > fatalLevel {
		level = fatalLevel
		zl.Level(level)
	}

	return &LoggerAdapter{
		zl:    zl,
		level: mrlog.Level(level),
	}
}

func (l *LoggerAdapter) Level() mrlog.Level {
	return mrlog.Level(l.zl.GetLevel())
}

func (l *LoggerAdapter) WithContext(ctx context.Context) context.Context {
	return mrlog.WithContext(ctx, l)
}

func (l *LoggerAdapter) With() mrlog.LoggerContext {
	return &ContextAdapter{zc: l.zl.With()}
}

func (l *LoggerAdapter) Debug() mrlog.LoggerEvent {
	if l.level > mrlog.DebugLevel {
		return (*EventAdapter)(nil)
	}

	return &EventAdapter{ze: l.zl.Debug()}
}

func (l *LoggerAdapter) Info() mrlog.LoggerEvent {
	if l.level > mrlog.InfoLevel {
		return (*EventAdapter)(nil)
	}

	return &EventAdapter{ze: l.zl.Info()}
}

func (l *LoggerAdapter) Warn() mrlog.LoggerEvent {
	if l.level > mrlog.WarnLevel {
		return (*EventAdapter)(nil)
	}

	return &EventAdapter{ze: l.zl.Warn()}
}

func (l *LoggerAdapter) Error() mrlog.LoggerEvent {
	if l.level > mrlog.ErrorLevel {
		return (*EventAdapter)(nil)
	}

	return &EventAdapter{ze: l.zl.Error()}
}

func (l *LoggerAdapter) Fatal() mrlog.LoggerEvent {
	return &EventAdapter{ze: l.zl.Fatal()}
}

func (l *LoggerAdapter) Panic() mrlog.LoggerEvent {
	return &EventAdapter{ze: l.zl.Panic()}
}

func (l *LoggerAdapter) Trace() mrlog.LoggerEvent {
	if l.level > mrlog.TraceLevel {
		return (*EventAdapter)(nil)
	}

	return &EventAdapter{ze: l.zl.Trace()}
}

func (l *LoggerAdapter) Printf(format string, args ...any) {
	l.zl.Printf(format, args...)
}

func (l *LoggerAdapter) Emit(ctx context.Context, eventName string, object any) {
	l.zl.Log().Str("event", eventName).Any("object", object).Send()
}

func (l *LoggerAdapter) EmitWithSource(ctx context.Context, eventName, source string, object any) {
	l.zl.Log().Str("event", eventName).Str("source", source).Any("object", object).Send()
}
