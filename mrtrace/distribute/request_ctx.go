package distribute

import (
	"context"
)

type (
	ctxRequestKey struct{}
)

// WithRequestID - устанавливает в контекст указанные идентификаторы процесса.
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, ctxRequestKey{}, requestID)
}

// RequestID - возвращает из контекста указанные идентификаторы текущего процесса.
func RequestID(ctx context.Context) string {
	if value, ok := ctx.Value(ctxRequestKey{}).(string); ok {
		return value
	}

	return ""
}
