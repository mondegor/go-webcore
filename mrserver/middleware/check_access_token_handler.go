package middleware

import (
	"net/http"

	"github.com/mondegor/go-sysmess/errors"
	"github.com/mondegor/go-sysmess/mrlog"

	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/request"
)

// CheckAccessTokenHandler - промежуточный обработчик запрещает доступ к обработчику авторизованному пользователю.
func CheckAccessTokenHandler(logger mrlog.Logger, handlerName string) func(next mrserver.HttpHandlerFunc) mrserver.HttpHandlerFunc {
	return func(next mrserver.HttpHandlerFunc) mrserver.HttpHandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) error {
			logger.Debug(r.Context(), "CheckAccessTokenHandler", "handler", handlerName)

			if accessToken := request.AccessToken(r); accessToken != "" {
				return errors.ErrHttpAccessForbidden
			}

			return next(w, r)
		}
	}
}
