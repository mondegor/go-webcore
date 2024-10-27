package mrapp

import (
	"context"
)

type (
	ctxProcessKey struct{}
)

// WithProcessContext - устанавливает в контекст указанный идентификатор процесса.
func WithProcessContext(ctx context.Context, processID string) context.Context {
	return context.WithValue(ctx, ctxProcessKey{}, processID)
}

// ProcessCtx - возвращает контекст с указанным идентификатором текущего процесса.
func ProcessCtx(ctx context.Context) string {
	if value, ok := ctx.Value(ctxProcessKey{}).(string); ok {
		return value
	}

	return ""
}
