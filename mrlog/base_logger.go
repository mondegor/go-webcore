package mrlog

import (
	"context"
	"log"
	"os"
)

type (
	// BaseLogger - логгер на крайний случай, например,
	// когда не был установлен логгер в контексте.
	BaseLogger struct {
		native  *log.Logger
		level   Level
		context []byte
	}
)

// Make sure the Image conforms with the mrlog.Logger interface.
var _ Logger = (*BaseLogger)(nil)

// New - создаёт объект BaseLogger.
func New(level Level) *BaseLogger {
	return &BaseLogger{
		native: log.Default(),
		level:  level,
	}
}

// Level - comment method.
func (l *BaseLogger) Level() Level {
	return l.level
}

// WithContext - comment method.
func (l *BaseLogger) WithContext(ctx context.Context) context.Context {
	return WithContext(ctx, l)
}

// With - comment method.
func (l *BaseLogger) With() LoggerContext {
	return &contextAdapter{
		logger: *l,
	}
}

// Info - comment method.
func (l *BaseLogger) Info() LoggerEvent {
	return l.newEvent(InfoLevel, nil)
}

// Debug - comment method.
func (l *BaseLogger) Debug() LoggerEvent {
	return l.newEvent(DebugLevel, nil)
}

// Warn - comment method.
func (l *BaseLogger) Warn() LoggerEvent {
	return l.newEvent(WarnLevel, nil)
}

// Error - comment method.
func (l *BaseLogger) Error() LoggerEvent {
	return l.newEvent(ErrorLevel, nil)
}

// Fatal - comment method.
func (l *BaseLogger) Fatal() LoggerEvent {
	return l.newEvent(FatalLevel, func(_ string) { os.Exit(1) }) //nolint:revive
}

// Panic - comment method.
func (l *BaseLogger) Panic() LoggerEvent {
	return l.newEvent(FatalLevel, func(message string) { panic(message) })
}

// Trace - comment method.
func (l *BaseLogger) Trace() LoggerEvent {
	return l.newEvent(TraceLevel, nil)
}

// Printf - comment method.
func (l *BaseLogger) Printf(format string, args ...any) {
	l.native.Printf(format+string(l.context), args...)
}

// newEvent - comment method.
func (l *BaseLogger) newEvent(level Level, done func(message string)) *eventAdapter {
	if l.level > level {
		return nil
	}

	return &eventAdapter{
		logger: l.native,
		buf:    l.context,
		done:   done,
	}
}
