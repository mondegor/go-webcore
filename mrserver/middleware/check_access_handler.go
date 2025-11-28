package middleware

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/mondegor/go-sysmess/mrerr/mr"
	"github.com/mondegor/go-sysmess/mrlog"

	"github.com/mondegor/go-webcore/mraccess"
	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/request"
)

// CheckAccessHandler - промежуточный обработчик проверки доступа к секции и конечному обработчику.
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
				return mr.ErrHttpClientUnauthorized.New()
			}

			currentUser, err := userProvider.UserByToken(ctx, accessToken)
			if err != nil {
				if mr.ErrUseCaseAccessForbidden.Is(err) {
					return mr.ErrHttpAccessForbidden.New()
				}

				return err
			}

			logger.Debug(ctx, "current user", "userId", uuid.UUID(currentUser.ID()).String())

			if action.Privilege != mraccess.PrivilegePublic && !currentUser.HasPrivilege(action.Privilege) {
				return mr.ErrHttpAccessForbidden.New()
			}

			if !currentUser.HasPermission(action.Permission) {
				return mr.ErrHttpAccessForbidden.New()
			}

			// замена языка переданного клиентом в заголовке Accept-Language
			// на язык, который был установлен пользователем
			if code := currentUser.LangCode(); code != "" {
				r.Header.Set(mrserver.HeaderKeyAcceptLanguage, code)
			}

			r.Header.Set(mrserver.HeaderKeyUserIDSlashGroup, uuid.UUID(currentUser.ID()).String()+"/"+currentUser.Group()) // userId/realm/kind

			if err = next(w, r); err != nil {
				if mr.ErrUseCaseAccessForbidden.Is(err) {
					return mr.ErrHttpAccessForbidden.New()
				}

				// если ошибка обработчика не связана с доступом к ресурсу
				return err
			}

			return nil
		}
	}
}
