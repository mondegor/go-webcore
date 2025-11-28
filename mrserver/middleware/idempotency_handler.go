package middleware

import (
	"net/http"

	"github.com/mondegor/go-sysmess/mrlog"

	"github.com/mondegor/go-webcore/mridempotency"
	"github.com/mondegor/go-webcore/mrserver"
)

// IdempotencyHandler - промежуточный обработчик для организации идемпотентных запросов.
func IdempotencyHandler(
	logger mrlog.Logger,
	provider mridempotency.Provider,
	sender mrserver.ResponseSender,
) func(next mrserver.HttpHandlerFunc) mrserver.HttpHandlerFunc {
	return func(next mrserver.HttpHandlerFunc) mrserver.HttpHandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) error {
			idempotencyKey := r.Header.Get(mrserver.HeaderKeyIdempotencyKey)

			if idempotencyKey == "" {
				return next(w, r)
			}

			if err := provider.Validate(idempotencyKey); err != nil {
				return err
			}

			cachedResponse, err := provider.Get(r.Context(), idempotencyKey)
			if err != nil {
				return err
			}

			if cachedResponse != nil {
				return sender.SendBytes(
					w,
					cachedResponse.StatusCode(),
					cachedResponse.Content(),
				)
			}

			unlock, err := provider.Lock(r.Context(), idempotencyKey)
			if err != nil {
				return err
			}

			defer unlock()

			sw := mrserver.NewCacheableResponseWriter(w)

			if err = next(sw, r); err != nil {
				return err
			}

			if err = provider.Store(r.Context(), idempotencyKey, sw); err != nil {
				logger.Error(r.Context(), "IdempotencyHandler->Store", "error", err)
			}

			return nil
		}
	}
}
