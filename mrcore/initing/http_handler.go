package initing

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/mondegor/go-core/errors"
	"github.com/mondegor/go-core/mraccess"
	"github.com/mondegor/go-core/mrlog"

	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/middleware"
	"github.com/mondegor/go-webcore/mrserver/request"
)

type (
	// PrepareHandlerFunc - функция-преобразователь HTTP-обработчика.
	// Используется для модификации обработчиков: установка разрешений,
	// добавление middleware, изменение URL и т.д.
	PrepareHandlerFunc func(handler mrserver.HttpHandler) mrserver.HttpHandler
)

// ==================|====================|============================|=================|===============|
// Group privilege   | Handler permission |           Actions          | Access errors   | Init/set user |
// ==================|====================|============================|=================|===============|
//  public           | everyone           | clear internal headers     | no              | no            |
// ------------------|--------------------|----------------------------|-----------------|---------------|
//  public           | guest-only         | check token                | if exists: 403  | no            |
// ------------------|--------------------|----------------------------|-----------------|---------------|
//  public           | any-user           | check token (auth only)    | 401             | yes           |
// ------------------|--------------------|----------------------------|-----------------|---------------|
//  public           | {permission}       | check token/permission     | 401, 403        | yes           |
// ------------------|--------------------|----------------------------|-----------------|---------------|
//  {privilege}      | everyone           | check token/priv           | 401, 403        | yes           |
// ------------------|--------------------|----------------------------|-----------------|---------------|
//  {privilege}      | guest-only         | warning, skip              | 403             |      ---      |
// ------------------|--------------------|----------------------------|-----------------|---------------|
//  {privilege}      | any-user           | check token/priv (auth)    | 401, 403        | yes           |
// ------------------|--------------------|----------------------------|-----------------|---------------|
//  {privilege}      | {permission}       | check token/priv/perm      | 401, 403        | yes           |
// ==================|====================|============================|=================|===============|

// WithPermission - создаёт функцию-преобразователь, которая устанавливает обработчику
// указанное разрешение (permission), если оно ещё не установлено.
func WithPermission(permission string) PrepareHandlerFunc {
	return func(handler mrserver.HttpHandler) mrserver.HttpHandler {
		if handler.Permission == "" {
			handler.Permission = permission
		}

		return handler
	}
}

// WithCheckAccessMiddleware - создаёт функцию-преобразователь, которая добавляет к обработчику
// middleware проверки доступа. Middleware проверяет токен доступа, привилегии и разрешения пользователя.
//
// Логика работы зависит от комбинации Privilege группы и Permission обработчика:
//   - public + everyone: без проверок (доступно всем, включая гостя), из запроса удаляются внутренние заголовки;
//   - privilege + everyone: доступ по наличию привилегии (401 без токена, 403 без привилегии);
//   - public + guest-only: проверка токена (возврат 403 если токен существует), из запроса удаляются внутренние заголовки;
//   - any-user: требуется только аутентификация (возврат 401 без токена), в приватной группе ещё и привилегия (403);
//   - public/privilege + permission: проверка токена/привилегии/разрешения (возврат 401/403);
//   - privilege + guest-only: предупреждение в лог и возврат 403.
//
// Служебные разрешения everyone и any-user передаются в systemPermissions провайдера ролей,
// поэтому они есть у каждой роли. Незарегистрированное разрешение (в том числе служебное,
// если его забыли передать) приводит к предупреждению в лог и возврату 403.
//
// Если userProvider не задан, в лог пишется ошибка, а обработчики, которым требуется пользователь,
// всегда возвращают 401 (public + everyone/guest-only работают как обычно, а privilege + guest-only
// и незарегистрированное разрешение по-прежнему дают предупреждение в лог и 403).
func WithCheckAccessMiddleware(
	logger mrlog.Logger,
	actionGroup mraccess.ActionGroup,
	userProvider mraccess.UserProvider,
	rightsAvailability mraccess.RightsRegistry,
) PrepareHandlerFunc {
	if actionGroup.Privilege != mraccess.PrivilegePublic && !rightsAvailability.IsRegistered(actionGroup.Privilege) {
		mrlog.Warn(
			logger,
			fmt.Sprintf(
				"Privilege '%s' is not registered for actionGroup '%s', perhaps, it is not registered in the config or is not associated with any role",
				actionGroup.Privilege, actionGroup.Name,
			),
		)
	}

	if userProvider == nil {
		mrlog.Error(
			logger,
			"UserProvider is not set for actionGroup",
			"actionGroup", actionGroup.Name,
			"error", errors.ErrInternalNilPointer.New(),
		)
	}

	return func(handler mrserver.HttpHandler) mrserver.HttpHandler {
		handler.URL = actionGroup.BasePath + strings.TrimLeft(handler.URL, "/")

		// everyone: доступно всем без проверок
		if actionGroup.Privilege == mraccess.PrivilegePublic && handler.Permission == mraccess.PermissionEveryone {
			handler.Func = middleware.ClearInternalHeadersHandler()(handler.Func)

			return handler
		}

		if actionGroup.Privilege == mraccess.PrivilegePublic && handler.Permission == mraccess.PermissionGuestOnly {
			next := middleware.ClearInternalHeadersHandler()(handler.Func)

			handler.Func = func(w http.ResponseWriter, r *http.Request) error {
				// guest-only доступен только неавторизованным пользователям
				if request.AccessToken(r) != "" {
					return errors.ErrHttpAccessForbidden
				}

				return next(w, r)
			}

			return handler
		}

		// guest-only не имеет смысла в приватной группе
		if actionGroup.Privilege != mraccess.PrivilegePublic && handler.Permission == mraccess.PermissionGuestOnly {
			mrlog.Warn(
				logger,
				"This permission cannot be present in the private actionGroup",
				"permission", handler.Permission,
				"method", handler.Method,
				"url", handler.URL,
			)

			handler.Func = func(_ http.ResponseWriter, _ *http.Request) error {
				return errors.ErrHttpAccessForbidden
			}

			return handler
		}

		if !rightsAvailability.IsRegistered(handler.Permission) {
			mrlog.Warn(
				logger,
				"Permission is not registered, perhaps, it is not registered in the config or is not associated with any role",
				"permission", handler.Permission,
				"method", handler.Method,
				"url", handler.URL,
			)

			handler.Func = func(_ http.ResponseWriter, _ *http.Request) error {
				return errors.ErrHttpAccessForbidden
			}

			return handler
		}

		// без UserProvider пользователя невозможно аутентифицировать
		if userProvider == nil {
			handler.Func = func(_ http.ResponseWriter, _ *http.Request) error {
				return errors.ErrHttpClientUnauthorized
			}

			return handler
		}

		handler.Func = middleware.CheckAccessHandler(
			logger,
			mraccess.Action{
				Name:       handler.URL,
				Privilege:  actionGroup.Privilege,
				Permission: handler.Permission,
			},
			userProvider,
		)(handler.Func)

		return handler
	}
}
