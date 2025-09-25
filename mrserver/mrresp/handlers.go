package mrresp

import (
	"context"
	"net/http"

	"github.com/mondegor/go-sysmess/mrlog"

	"github.com/mondegor/go-webcore/mrlib"
)

// HandlerGetStatusOkAsJSON - возвращает обработчик, который формирует ответ OK в JSON формате.
func HandlerGetStatusOkAsJSON(logger mrlog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		mrlib.Write(r.Context(), logger, w, []byte("{\"status\": \"OK\"}"))
	}
}

// HandlerGetHealth - обработчик для использования в качестве проверки работоспособности сервиса.
func HandlerGetHealth(isAvailable func(ctx context.Context) bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status := http.StatusOK

		if !isAvailable(r.Context()) {
			status = http.StatusServiceUnavailable
		}

		w.WriteHeader(status)
	}
}
