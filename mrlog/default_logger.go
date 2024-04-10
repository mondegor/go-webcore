package mrlog

import (
	"context"
	"log"
	"os"
)

const (
	defaultLoggerPrefix = "[mr-default-logger] "
)

type (
	DefaultLogger struct {
		native  *log.Logger
		level   Level
		context []byte
	}
)

func New(level Level) *DefaultLogger {
	return &DefaultLogger{
		native: log.Default(),
		level:  level,
	}
}

func (l *DefaultLogger) Level() Level {
	return l.level
}

func (l *DefaultLogger) WithContext(ctx context.Context) context.Context {
	return WithContext(ctx, l)
}

func (l *DefaultLogger) With() LoggerContext {
	return &defaultContext{
		logger: *l,
	}
}

func (l *DefaultLogger) Info() LoggerEvent {
	return l.newEvent(InfoLevel, nil)
}

func (l *DefaultLogger) Debug() LoggerEvent {
	return l.newEvent(DebugLevel, nil)
}

func (l *DefaultLogger) Warn() LoggerEvent {
	return l.newEvent(WarnLevel, nil)
}

func (l *DefaultLogger) Error() LoggerEvent {
	return l.newEvent(ErrorLevel, nil)
}

func (l *DefaultLogger) Fatal() LoggerEvent {
	return l.newEvent(FatalLevel, func(message string) { os.Exit(1) })
}

func (l *DefaultLogger) Panic() LoggerEvent {
	return l.newEvent(FatalLevel, func(message string) { panic(message) })
}

func (l *DefaultLogger) Trace() LoggerEvent {
	return l.newEvent(TraceLevel, nil)
}

func (l *DefaultLogger) Printf(format string, args ...any) {
	l.native.Printf(defaultLoggerPrefix+format+string(l.context), args...)
}

func (l *DefaultLogger) Emit(ctx context.Context, eventName string, object any) {
	l.native.Printf("event=%s, object=%v", eventName, object)
}

func (l *DefaultLogger) EmitWithSource(ctx context.Context, eventName, source string, object any) {
	l.native.Printf("event=%s, source=%s, object=%v", eventName, source, object)
}

func (l *DefaultLogger) newEvent(level Level, done func(message string)) *defaultEvent {
	if l.level > level {
		return nil
	}

	return &defaultEvent{
		logger: l.native,
		buf:    l.context,
		done:   done,
	}
}
