package middleware

import (
	"net/http"

	"github.com/mondegor/go-webcore/mrserver"
)

// ClearInternalHeadersHandler - middleware, удаляющий из запроса служебные заголовки (X-Internal-*).
// Используется для эндпоинтов, доступных всем без проверки доступа: внутренние
// заголовки устанавливает только сервер, поэтому значения, подставленные клиентом,
// не должны дойти до обработчика.
func ClearInternalHeadersHandler() func(next mrserver.HttpHandlerFunc) mrserver.HttpHandlerFunc {
	return func(next mrserver.HttpHandlerFunc) mrserver.HttpHandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) error {
			r.Header.Del(mrserver.HeaderKeyInternalUserIDSlashGroup)
			r.Header.Del(mrserver.HeaderKeyInternalSessionID)
			r.Header.Del(mrserver.HeaderKeyInternalLangCode)
			r.Header.Del(mrserver.HeaderKeyInternalTimeZone)

			return next(w, r)
		}
	}
}
