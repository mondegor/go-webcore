package zerolog

import (
	"context"

	"github.com/rs/zerolog"

	"github.com/mondegor/go-webcore/mrlog"
)

// go get -u github.com/rs/zerolog

type (
	// LoggerAdapter - comment struct.
	LoggerAdapter struct {
		zl   zerolog.Logger
		opts loggerOptions
	}

	// Options - comment struct.
	Options struct {
		Zerolog          zerolog.Logger
		Level            mrlog.Level
		PrepareErrorFunc func(err error) error
	}

	loggerOptions struct {
		level            mrlog.Level
		prepareErrorFunc func(err error) error
	}
)

// New - создаёт объект LoggerAdapter.
func New(opts Options) *LoggerAdapter {
	if opts.PrepareErrorFunc == nil {
		opts.PrepareErrorFunc = func(err error) error {
			return err
		}
	}

	return &LoggerAdapter{
		zl: opts.Zerolog.Level(zerolog.Level(opts.Level)),
		opts: loggerOptions{
			level:            opts.Level,
			prepareErrorFunc: opts.PrepareErrorFunc,
		},
	}
}

// Level - comment method.
func (l *LoggerAdapter) Level() mrlog.Level {
	return l.opts.level
}

// WithContext - comment method.
func (l *LoggerAdapter) WithContext(ctx context.Context) context.Context {
	return mrlog.WithContext(ctx, l)
}

// With - comment method.
func (l *LoggerAdapter) With() mrlog.LoggerContext {
	return &contextAdapter{
		zc:   l.zl.With(),
		opts: l.opts,
	}
}

// Debug - comment method.
func (l *LoggerAdapter) Debug() mrlog.LoggerEvent {
	if l.opts.level > mrlog.DebugLevel {
		return (*eventAdapter)(nil)
	}

	return l.newEventAdapter(l.zl.Debug())
}

// Info - comment method.
func (l *LoggerAdapter) Info() mrlog.LoggerEvent {
	if l.opts.level > mrlog.InfoLevel {
		return (*eventAdapter)(nil)
	}

	return l.newEventAdapter(l.zl.Info())
}

// Warn - comment method.
func (l *LoggerAdapter) Warn() mrlog.LoggerEvent {
	if l.opts.level > mrlog.WarnLevel {
		return (*eventAdapter)(nil)
	}

	return l.newEventAdapter(l.zl.Warn())
}

// Error - comment method.
func (l *LoggerAdapter) Error() mrlog.LoggerEvent {
	if l.opts.level > mrlog.ErrorLevel {
		return (*eventAdapter)(nil)
	}

	return l.newEventAdapter(l.zl.Error())
}

// Fatal - comment method.
func (l *LoggerAdapter) Fatal() mrlog.LoggerEvent {
	return l.newEventAdapter(l.zl.Fatal())
}

// Panic - comment method.
func (l *LoggerAdapter) Panic() mrlog.LoggerEvent {
	return l.newEventAdapter(l.zl.Panic())
}

// Trace - comment method.
func (l *LoggerAdapter) Trace() mrlog.LoggerEvent {
	if l.opts.level > mrlog.TraceLevel {
		return (*eventAdapter)(nil)
	}

	return l.newEventAdapter(l.zl.Trace())
}

// Printf - comment method.
func (l *LoggerAdapter) Printf(format string, args ...any) {
	l.zl.Printf(format, args...)
}

// newEventAdapter - comment method.
func (l *LoggerAdapter) newEventAdapter(ze *zerolog.Event) *eventAdapter {
	return &eventAdapter{
		ze:               ze,
		prepareErrorFunc: l.opts.prepareErrorFunc,
	}
}
