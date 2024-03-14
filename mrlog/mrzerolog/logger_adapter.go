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
		zl   zerolog.Logger
		opts loggerOptions
	}

	loggerOptions struct {
		level              mrlog.Level
		isAutoCallerOnFunc func(err error) bool
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

	if opts.IsAutoCallerOnFunc == nil {
		opts.IsAutoCallerOnFunc = func(err error) bool {
			return true
		}
	}

	return &LoggerAdapter{
		zl: logger,
		opts: loggerOptions{
			level:              opts.Level,
			isAutoCallerOnFunc: opts.IsAutoCallerOnFunc,
		},
	}
}

func (l *LoggerAdapter) Level() mrlog.Level {
	return l.opts.level
}

func (l *LoggerAdapter) WithContext(ctx context.Context) context.Context {
	return mrlog.WithContext(ctx, l)
}

func (l *LoggerAdapter) With() mrlog.LoggerContext {
	return &ContextAdapter{zc: l.zl.With(), opts: l.opts}
}

func (l *LoggerAdapter) Debug() mrlog.LoggerEvent {
	if l.opts.level > mrlog.DebugLevel {
		return (*EventAdapter)(nil)
	}

	return l.newEventAdapter(l.zl.Debug(), false)
}

func (l *LoggerAdapter) Info() mrlog.LoggerEvent {
	if l.opts.level > mrlog.InfoLevel {
		return (*EventAdapter)(nil)
	}

	return l.newEventAdapter(l.zl.Info(), false)
}

func (l *LoggerAdapter) Warn() mrlog.LoggerEvent {
	if l.opts.level > mrlog.WarnLevel {
		return (*EventAdapter)(nil)
	}

	return l.newEventAdapter(l.zl.Warn(), false)
}

func (l *LoggerAdapter) Error() mrlog.LoggerEvent {
	if l.opts.level > mrlog.ErrorLevel {
		return (*EventAdapter)(nil)
	}

	return l.newEventAdapter(l.zl.Error(), true)
}

func (l *LoggerAdapter) Fatal() mrlog.LoggerEvent {
	return l.newEventAdapter(l.zl.Fatal(), true)
}

func (l *LoggerAdapter) Panic() mrlog.LoggerEvent {
	return l.newEventAdapter(l.zl.Panic(), false)
}

func (l *LoggerAdapter) Trace() mrlog.LoggerEvent {
	if l.opts.level > mrlog.TraceLevel {
		return (*EventAdapter)(nil)
	}

	return l.newEventAdapter(l.zl.Trace(), false)
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

func (l *LoggerAdapter) newEventAdapter(ze *zerolog.Event, isAutoCallerAllowed bool) *EventAdapter {
	return &EventAdapter{
		ze:                  ze,
		isAutoCallerAllowed: isAutoCallerAllowed,
		isAutoCallerOnFunc:  l.opts.isAutoCallerOnFunc,
	}
}
