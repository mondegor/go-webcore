package distribute

import (
	"context"
)

type (
	ctxTaskKey struct{}
)

// WithTaskID - устанавливает в контекст указанные идентификаторы процесса.
func WithTaskID(ctx context.Context, taskID string) context.Context {
	return context.WithValue(ctx, ctxTaskKey{}, taskID)
}

// TaskID - возвращает из контекста указанные идентификаторы текущего процесса.
func TaskID(ctx context.Context) string {
	if value, ok := ctx.Value(ctxTaskKey{}).(string); ok {
		return value
	}

	return ""
}
