package sentry

import (
	"context"
	"log/slog"
)

// Wrapper - обёртка над slog.Handler для интеграции с Sentry.
// Позволяет перехватывать лог-записи перед их отправкой в Sentry.
// Реализует интерфейс slog.Handler.
type (
	Wrapper struct {
		handler slog.Handler
	}
)

// NewWrapper - создаёт обёртку Wrapper для указанного slog.Handler.
func NewWrapper(handler slog.Handler) *Wrapper {
	return &Wrapper{
		handler: handler,
	}
}

// Enabled - проверяет, включён ли обработчик для заданного уровня логирования.
func (w *Wrapper) Enabled(ctx context.Context, level slog.Level) bool {
	return w.handler.Enabled(ctx, level)
}

// Handle - обрабатывает лог-запись, проходя по всем её атрибутам.
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

// WithAttrs - создаёт новый обработчик с дополнительными атрибутами.
func (w *Wrapper) WithAttrs(attrs []slog.Attr) slog.Handler {
	return w.handler.WithAttrs(attrs)
}

// WithGroup - создаёт новый обработчик с указанной группой атрибутов.
func (w *Wrapper) WithGroup(name string) slog.Handler {
	return w.handler.WithGroup(name)
}
