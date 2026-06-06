package middleware

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/mondegor/go-sysmess/errors"
	"github.com/mondegor/go-sysmess/mraccess"
	"github.com/mondegor/go-sysmess/mrlog"

	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/request"
)

// CheckAccessHandler - middleware для проверки доступа пользователя к конечному обработчику запроса.
//
// Логика работы:
//  1. Извлекает access token из запроса;
//  2. Получает данные пользователя через userProvider;
//  3. Проверяет привилегию (Privilege) и разрешение (Permission) пользователя;
//  4. Устанавливает заголовки: Accept-Language (из профиля пользователя), UserID/Group;
//  5. Вызывает следующий обработчик в цепочке.
func CheckAccessHandler(
	logger mrlog.Logger,
	action mraccess.Action,
	userProvider mraccess.UserProvider,
) func(next mrserver.HttpHandlerFunc) mrserver.HttpHandlerFunc {
	return func(next mrserver.HttpHandlerFunc) mrserver.HttpHandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) error {
			ctx := r.Context()
			logger.Debug(
				ctx,
				"CheckAccessHandler",
				"handler", action.Name,
				"privilege", action.Privilege,
				"permission", action.Permission,
			)

			accessToken := request.AccessToken(r)
			if accessToken == "" {
				return errors.ErrHttpClientUnauthorized
			}

			currentUser, err := userProvider.UserByToken(ctx, accessToken)
			if err != nil {
				// if errors.Is(err, errors.ErrAccessForbidden) {
				// 	return errors.ErrHttpAccessForbidden
				// }
				return err
			}

			logger.Debug(ctx, "current user", "userId", uuid.UUID(currentUser.ID()).String())

			if action.Privilege != mraccess.PrivilegePublic && !currentUser.HasPrivilege(action.Privilege) {
				return errors.ErrHttpAccessForbidden
			}

			if !currentUser.HasPermission(action.Permission) {
				return errors.ErrHttpAccessForbidden
			}

			// замена языка переданного клиентом в заголовке Accept-Language
			// на язык, который был установлен пользователем
			if code := currentUser.LangCode(); code != "" {
				r.Header.Set(mrserver.HeaderKeyAcceptLanguage, code)
			}

			r.Header.Set(mrserver.HeaderKeyUserIDSlashGroup, uuid.UUID(currentUser.ID()).String()+"/"+currentUser.Group()) // userId/realm/kind

			if err = next(w, r); err != nil {
				// if errors.Is(err, errors.ErrAccessForbidden) {
				// 	return errors.ErrHttpAccessForbidden
				// }
				// если ошибка обработчика не связана с доступом к ресурсу
				return err
			}

			return nil
		}
	}
}
