package noplog

import (
	"context"
	"os"

	"github.com/mondegor/go-webcore/mrlog"
)

type (
	// LoggerAdapter - comment struct.
	LoggerAdapter struct{}
)

// Make sure the Image conforms with the mrlog.Logger interface.
var _ mrlog.Logger = (*LoggerAdapter)(nil)

// New - создаёт объект LoggerAdapter.
func New() *LoggerAdapter {
	return &LoggerAdapter{}
}

// Level - comment method.
func (l *LoggerAdapter) Level() mrlog.Level {
	return mrlog.FatalLevel
}

// WithContext - comment method.
func (l *LoggerAdapter) WithContext(ctx context.Context) context.Context {
	return mrlog.WithContext(ctx, l)
}

// With - comment method.
func (l *LoggerAdapter) With() mrlog.LoggerContext {
	return &contextAdapter{
		logger: *l,
	}
}

// Info - comment method.
func (l *LoggerAdapter) Info() mrlog.LoggerEvent {
	return (*eventAdapter)(nil)
}

// Debug - comment method.
func (l *LoggerAdapter) Debug() mrlog.LoggerEvent {
	return (*eventAdapter)(nil)
}

// Warn - comment method.
func (l *LoggerAdapter) Warn() mrlog.LoggerEvent {
	return (*eventAdapter)(nil)
}

// Error - comment method.
func (l *LoggerAdapter) Error() mrlog.LoggerEvent {
	return (*eventAdapter)(nil)
}

// Fatal - comment method.
func (l *LoggerAdapter) Fatal() mrlog.LoggerEvent {
	return &eventAdapter{
		done: func() {
			os.Exit(1) //nolint:revive
		},
	}
}

// Panic - comment method.
func (l *LoggerAdapter) Panic() mrlog.LoggerEvent {
	return &eventAdapter{
		done: func() {
			panic("noplog.LoggerAdapter")
		},
	}
}

// Trace - comment method.
func (l *LoggerAdapter) Trace() mrlog.LoggerEvent {
	return (*eventAdapter)(nil)
}

// Printf - comment method.
func (l *LoggerAdapter) Printf(_ string, _ ...any) {
}
