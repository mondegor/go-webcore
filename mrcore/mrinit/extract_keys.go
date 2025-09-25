package mrinit

import (
	"context"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrtrace/distribute"
)

// ExtractKeys - возвращает попарно (processKey/processId)
// все имеющиеся ID процессов из указанного контекста.
func ExtractKeys(ctx context.Context) (keyValue []any) {
	if ctx == nil || ctx == context.Background() {
		return nil
	}

	keyValue = make([]any, 0, 4)

	if value := distribute.CorrelationID(ctx); value != "" {
		keyValue = append(keyValue, mrcore.KeyCorrelationID, value)
	}

	if value := distribute.ProcessID(ctx); value != "" {
		keyValue = append(keyValue, mrcore.KeyProcessID, value)
	}

	if value := distribute.RequestID(ctx); value != "" {
		return append(keyValue, mrcore.KeyRequestID, value) // если задан RequestID, то не должно быть WorkerID или TaskID
	}

	if value := distribute.WorkerID(ctx); value != "" {
		keyValue = append(keyValue, mrcore.KeyWorkerID, value)
	}

	if value := distribute.TaskID(ctx); value != "" {
		keyValue = append(keyValue, mrcore.KeyTaskID, value)
	}

	return keyValue
}
