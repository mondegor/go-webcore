package middleware

import (
	"net/http"
	"time"

	"github.com/mondegor/go-sysmess/mrlog"

	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/observe"
)

const (
	// :TODO: вынести в настройки.
	traceRequestBodyMaxLen  = 2048
	traceResponseBodyMaxLen = 2048
)

// ObserverHandler - промежуточный обработчик, который собирает статистику запросов.
func ObserverHandler(
	logger mrlog.Logger,
	observer mrserver.RequestStat,
) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			logger.Info(r.Context(), "REQUEST",
				"method", r.Method,
				"uri", r.RequestURI,
			)

			sr := observe.NewRequestReader(r, traceRequestBodyMaxLen)
			sw := observe.NewResponseWriter(w, traceResponseBodyMaxLen)

			defer func() {
				observer.Emit(
					sr.Request(),
					sr.Content(),
					sr.Size(),
					sw.Content(),
					sw.Size(),
					time.Since(start),
					sw.StatusCode(),
				)
			}()

			next.ServeHTTP(sw, sr.Request())
		})
	}
}
