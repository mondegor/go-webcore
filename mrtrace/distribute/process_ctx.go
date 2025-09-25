package distribute

import (
	"context"
)

type (
	ctxProcessKey struct{}
)

// WithProcessID - устанавливает в контекст указанные идентификаторы процесса.
func WithProcessID(ctx context.Context, processID string) context.Context {
	return context.WithValue(ctx, ctxProcessKey{}, processID)
}

// ProcessID - возвращает из контекста указанные идентификаторы текущего процесса.
func ProcessID(ctx context.Context) string {
	if value, ok := ctx.Value(ctxProcessKey{}).(string); ok {
		return value
	}

	return ""
}
