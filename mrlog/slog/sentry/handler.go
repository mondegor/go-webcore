package sentry

import (
	"context"
	"log/slog"
)

// Wrapper - обёртка реализующая интерфейс slog.Handler.
type (
	Wrapper struct {
		handler slog.Handler
	}
)

// NewWrapper - создаёт объект Wrapper.
func NewWrapper(handler slog.Handler) *Wrapper {
	return &Wrapper{
		handler: handler,
	}
}

// Enabled - comment func.
func (w *Wrapper) Enabled(ctx context.Context, level slog.Level) bool {
	return w.handler.Enabled(ctx, level)
}

// Handle - comment func.
func (w *Wrapper) Handle(ctx context.Context, record slog.Record) error {
	record.Attrs(func(_ slog.Attr) bool {
		// if attr.Value.Kind() == slog.KindAny {
		// 	if err, ok := attr.Value.Any().(error); ok {
		// 		record.
		// 	}
		// }
		return true
	})

	return w.handler.Handle(ctx, record)
}

// WithAttrs - comment func.
func (w *Wrapper) WithAttrs(attrs []slog.Attr) slog.Handler {
	return w.handler.WithAttrs(attrs)
}

// WithGroup - comment func.
func (w *Wrapper) WithGroup(name string) slog.Handler {
	return w.handler.WithGroup(name)
}
