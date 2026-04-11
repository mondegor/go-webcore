package middleware

import (
	"net/http"

	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-sysmess/mrtrace"

	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/request"
)

// RequestIDHandler - middleware для генерации и управления идентификаторами запросов.
//
// Логика работы:
//  1. Генерирует уникальный RequestID и добавляет его в контекст;
//  2. Устанавливает RequestID в заголовок ответа (X-Request-ID);
//  3. Извлекает CorrelationID из запроса (если присутствует);
//  4. Добавляет CorrelationID в контекст и заголовок ответа;
//  5. Передаёт управление следующему обработчику с обновлённым контекстом.
//
// Важно:
//   - RequestID генерируется автоматически для каждого запроса;
//   - CorrelationID опционален и берётся из входящего запроса.
func RequestIDHandler(logger mrlog.Logger, traceManager mrtrace.ContextManager) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := traceManager.WithGeneratedProcessID(r.Context(), mrtrace.KeyRequestID)
			w.Header().Set(mrserver.HeaderKeyRequestID, traceManager.ProcessID(ctx, mrtrace.KeyRequestID))

			if correlationID, err := request.CorrelationID(r); err != nil {
				logger.Warn(ctx, "RequestIDHandler", "error", err)
			} else if correlationID != "" {
				ctx = traceManager.WithProcessID(ctx, mrtrace.KeyCorrelationID, correlationID)
				w.Header().Set(mrserver.HeaderKeyCorrelationID, correlationID)
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
