package mrlogbase

import (
	"context"
	"log"
	"os"

	"github.com/mondegor/go-webcore/mrlog"
)

type (
	// LoggerAdapter - comment struct.
	LoggerAdapter struct {
		native  *log.Logger
		level   mrlog.Level
		context []byte
	}
)

// Make sure the Image conforms with the mrlog.Logger interface.
var _ mrlog.Logger = (*LoggerAdapter)(nil)

// New - создаёт объект LoggerAdapter.
func New(level mrlog.Level) *LoggerAdapter {
	return &LoggerAdapter{
		native: log.Default(),
		level:  level,
	}
}

// Level - comment method.
func (l *LoggerAdapter) Level() mrlog.Level {
	return l.level
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
	return l.newEvent(mrlog.InfoLevel, nil)
}

// Debug - comment method.
func (l *LoggerAdapter) Debug() mrlog.LoggerEvent {
	return l.newEvent(mrlog.DebugLevel, nil)
}

// Warn - comment method.
func (l *LoggerAdapter) Warn() mrlog.LoggerEvent {
	return l.newEvent(mrlog.WarnLevel, nil)
}

// Error - comment method.
func (l *LoggerAdapter) Error() mrlog.LoggerEvent {
	return l.newEvent(mrlog.ErrorLevel, nil)
}

// Fatal - comment method.
func (l *LoggerAdapter) Fatal() mrlog.LoggerEvent {
	return l.newEvent(mrlog.FatalLevel, func(_ string) { os.Exit(1) }) //nolint:revive
}

// Panic - comment method.
func (l *LoggerAdapter) Panic() mrlog.LoggerEvent {
	return l.newEvent(mrlog.FatalLevel, func(message string) { panic(message) })
}

// Trace - comment method.
func (l *LoggerAdapter) Trace() mrlog.LoggerEvent {
	return l.newEvent(mrlog.TraceLevel, nil)
}

// Printf - comment method.
func (l *LoggerAdapter) Printf(format string, args ...any) {
	l.native.Printf(format+string(l.context), args...)
}

// newEvent - comment method.
func (l *LoggerAdapter) newEvent(level mrlog.Level, done func(message string)) *eventAdapter {
	if l.level > level {
		return nil
	}

	return &eventAdapter{
		logger: l.native,
		buf:    l.context,
		done:   done,
	}
}
