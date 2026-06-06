package middleware

import (
	"net/http"

	"github.com/mondegor/go-sysmess/mridempotency"
	"github.com/mondegor/go-sysmess/mrlog"

	"github.com/mondegor/go-webcore/mrserver"
)

// IdempotencyHandler - middleware для обеспечения идемпотентности запросов.
//
// Логика работы:
//  1. Извлекает Idempotency-Key из заголовка запроса;
//  2. Если ключ отсутствует, пропускает запрос без обработки;
//  3. Валидирует ключ через provider.Validate();
//  4. Проверяет наличие кэшированного ответа через provider.Get();
//  5. Если ответ найден, возвращает его немедленно (без вызова next);
//  6. Если ответ не найден, блокирует ключ через provider.Lock();
//  7. Вызывает следующий обработчик, записывая ответ в CacheableResponseWriter;
//  8. Сохраняет ответ в кэш через provider.Save().
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

			if err = provider.Save(r.Context(), idempotencyKey, sw); err != nil {
				logger.Error(r.Context(), "IdempotencyHandler->Store", "error", err)
			}

			return nil
		}
	}
}
