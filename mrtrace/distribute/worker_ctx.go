package distribute

import (
	"context"
)

type (
	ctxWorkerKey struct{}
)

// WithWorkerID - устанавливает в контекст указанные идентификаторы процесса.
func WithWorkerID(ctx context.Context, workerID string) context.Context {
	return context.WithValue(ctx, ctxWorkerKey{}, workerID)
}

// WorkerID - возвращает из контекста указанные идентификаторы текущего процесса.
func WorkerID(ctx context.Context) string {
	if value, ok := ctx.Value(ctxWorkerKey{}).(string); ok {
		return value
	}

	return ""
}
