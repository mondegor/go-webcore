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

// Enabled - проверяет, включён ли обработчик для заданного уровня логирования.
func (w *Wrapper) Enabled(ctx context.Context, level slog.Level) bool {
	return w.handler.Enabled(ctx, level)
}

// Handle - обрабатывает запись журнала.
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

// WithAttrs - возвращает новый обработчик с добавлением указанных атрибутов.
func (w *Wrapper) WithAttrs(attrs []slog.Attr) slog.Handler {
	return w.handler.WithAttrs(attrs)
}

// WithGroup - возвращает новый обработчик с указанной группой атрибутов.
func (w *Wrapper) WithGroup(name string) slog.Handler {
	return w.handler.WithGroup(name)
}
