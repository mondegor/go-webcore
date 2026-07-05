package mrresp

import (
	"context"
	"net/http"

	"github.com/mondegor/go-core/mrlog"
	"github.com/mondegor/go-core/util/xio"
)

// HandlerGetStatusOkAsJSON - создаёт обработчик для ответов 200 OK.
// Возвращает JSON-ответ: {"status": "OK"}.
// Используется для простых проверок работоспособности (liveness probe).
func HandlerGetStatusOkAsJSON(logger mrlog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		xio.Write(r.Context(), logger, w, []byte("{\"status\": \"OK\"}"))
	}
}

// HandlerGetHealth - создаёт обработчик для проверки готовности сервиса (readiness probe).
// Возвращает 200 OK если сервис готов, или 503 Service Unavailable если нет.
func HandlerGetHealth(availableFunc func(ctx context.Context) bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status := http.StatusOK

		if !availableFunc(r.Context()) {
			status = http.StatusServiceUnavailable
		}

		w.WriteHeader(status)
	}
}
