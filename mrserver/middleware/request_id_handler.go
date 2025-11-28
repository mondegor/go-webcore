package middleware

import (
	"net/http"

	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-sysmess/mrtrace"

	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/request"
)

// RequestIDHandler - промежуточный обработчик,
// который устанавливает в контекст requestId, correlationId.
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

			// необходимо гарантировать, чтобы этот заголовок не был передан из вне
			r.Header.Del(mrserver.HeaderKeyUserIDSlashGroup)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
