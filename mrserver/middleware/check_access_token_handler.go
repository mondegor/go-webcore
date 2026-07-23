package middleware

import (
	"net/http"

	"github.com/mondegor/go-core/errors"
	"github.com/mondegor/go-core/mrlog"

	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/request"
)

// CheckAccessTokenHandler - middleware, запрещающий доступ авторизованным пользователям.
// Используется для эндпоинтов, доступных только неавторизованным пользователям
// (например: регистрация, вход в систему). Если запрос содержит access token,
// доступ будет запрещён.
func CheckAccessTokenHandler(logger mrlog.Logger, handlerName string) func(next mrserver.HttpHandlerFunc) mrserver.HttpHandlerFunc {
	return func(next mrserver.HttpHandlerFunc) mrserver.HttpHandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) error {
			logger.Debug(r.Context(), "CheckAccessTokenHandler", "handler", handlerName)

			if accessToken := request.AccessToken(r); accessToken != "" {
				return errors.ErrHttpAccessForbidden
			}

			// внутренние заголовки удаляются, чтобы клиент не мог выдать себя
			// за авторизованного пользователя, подставив их самостоятельно
			r.Header.Del(mrserver.HeaderKeyInternalUserIDSlashGroup)
			r.Header.Del(mrserver.HeaderKeyInternalSessionID)
			r.Header.Del(mrserver.HeaderKeyInternalLangCode)
			r.Header.Del(mrserver.HeaderKeyInternalTimeZone)

			return next(w, r)
		}
	}
}
