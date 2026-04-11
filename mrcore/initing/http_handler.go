package initing

import (
	"fmt"
	"net/http"

	"github.com/mondegor/go-sysmess/errors"
	"github.com/mondegor/go-sysmess/mrlog"

	"github.com/mondegor/go-webcore/mraccess"
	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/middleware"
)

type (
	// PrepareHandlerFunc - функция-преобразователь HTTP-обработчика.
	// Используется для модификации обработчиков: установка разрешений,
	// добавление middleware, изменение URL и т.д.
	PrepareHandlerFunc func(handler mrserver.HttpHandler) mrserver.HttpHandler
)

// ==================|====================|========================|=================|===============|
// Group privilege   | Handler permission |         Actions        | Access errors   | Init/set user |
// ==================|====================|========================|=================|===============|
//  public           | guest              | no                     | no              | no            |
// ------------------|--------------------|------------------------|-----------------|---------------|
//  public           | guest-only         | check token            | if exists: 403  | no            |
// ------------------|--------------------|------------------------|-----------------|---------------|
//  public           | {permission}       | check token/permission | 401, 403        | yes           |
// ------------------|--------------------|------------------------|-----------------|---------------|
//  {privilege}      | guest              | warning, skip          |       ---       |      ---      |
// ------------------|--------------------|------------------------|-----------------|---------------|
//  {privilege}      | guest-only         | warning, skip          |       ---       |      ---      |
// ------------------|--------------------|------------------------|-----------------|---------------|
//  {privilege}      | {permission}       | check token/priv/perm  | 401, 403        | yes           |
// ==================|====================|========================|=================|===============|

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
//   - public + guest: без проверок;
//   - public + guest-only: проверка токена (возврат 403 если токен существует);
//   - public + permission: проверка токена и разрешения (возврат 401/403);
//   - privilege + permission: проверка токена, привилегии и разрешения (возврат 401/403);
//   - privilege + guest/guest-only: предупреждение в лог и возврат 403.
func WithCheckAccessMiddleware(
	logger mrlog.Logger,
	actionGroup *mraccess.ActionGroup,
	userProvider mraccess.UserProvider,
	rightsAvailability mraccess.RightsChecker,
) PrepareHandlerFunc {
	if actionGroup.Privilege != mraccess.PrivilegePublic && !rightsAvailability.HasPrivilege(actionGroup.Privilege) {
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

		return func(handler mrserver.HttpHandler) mrserver.HttpHandler {
			handler.Func = func(_ http.ResponseWriter, _ *http.Request) error {
				return errors.ErrHttpClientUnauthorized
			}

			return handler
		}
	}

	return func(handler mrserver.HttpHandler) mrserver.HttpHandler {
		handler.URL = actionGroup.BasePath.BuildPath(handler.URL)

		if actionGroup.Privilege == mraccess.PrivilegePublic && handler.Permission == mraccess.PermissionAnyUser {
			return handler
		}

		if actionGroup.Privilege == mraccess.PrivilegePublic && handler.Permission == mraccess.PermissionGuestOnly {
			middlewareFunc := middleware.CheckAccessTokenHandler(
				logger,
				handler.URL,
			)

			handler.Func = middlewareFunc(handler.Func)

			return handler
		}

		if actionGroup.Privilege != mraccess.PrivilegePublic &&
			(handler.Permission == mraccess.PermissionAnyUser || handler.Permission == mraccess.PermissionGuestOnly) {
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

		if !rightsAvailability.HasPermission(handler.Permission) {
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
