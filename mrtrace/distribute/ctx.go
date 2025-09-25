package distribute

import (
	"context"
)

// NewContextWithIDs - возвращает новый чистый контекст
// со всеми ID процессов, скопированными из указанного контекста.
func NewContextWithIDs(originalCtx context.Context) context.Context {
	ctx := context.Background()

	if originalCtx == nil || originalCtx == ctx {
		return ctx
	}

	if value := CorrelationID(originalCtx); value != "" {
		ctx = WithCorrelationID(ctx, value)
	}

	if value := ProcessID(originalCtx); value != "" {
		ctx = WithProcessID(ctx, value)
	}

	if value := RequestID(originalCtx); value != "" {
		ctx = WithCorrelationID(ctx, value)

		return ctx // если задан RequestID, то не должно быть WorkerID или TaskID
	}

	if value := WorkerID(originalCtx); value != "" {
		ctx = WithWorkerID(ctx, value)
	}

	if value := TaskID(originalCtx); value != "" {
		ctx = WithTaskID(ctx, value)
	}

	return ctx
}

// FindCorrelationID - возвращает первый попавшийся ID из указанного контекста,
// который можно использовать в качестве CorrelationID.
func FindCorrelationID(ctx context.Context) string {
	if value := CorrelationID(ctx); value != "" {
		return value
	}

	if value := RequestID(ctx); value != "" {
		return value
	}

	if value := TaskID(ctx); value != "" {
		return value
	}

	if value := WorkerID(ctx); value != "" {
		return value
	}

	return ProcessID(ctx)
}
