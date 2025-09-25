package mrinit

import (
	"fmt"
	"net/http"

	"github.com/mondegor/go-sysmess/mrerr/mr"
	"github.com/mondegor/go-sysmess/mrlog"

	"github.com/mondegor/go-webcore/mraccess"
	"github.com/mondegor/go-webcore/mrserver"
)

// ==================|====================|========================|=================|===============|
// Section privilege | Handler permission |         Actions        | Access errors   | Init/set user |
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

// WithPermission - comment func.
func WithPermission(permission string) PrepareHandlerFunc {
	return func(handler mrserver.HttpHandler) mrserver.HttpHandler {
		if handler.Permission == "" {
			handler.Permission = permission
		}

		return handler
	}
}

// WithMiddlewareCheckAccess - comment func.
func WithMiddlewareCheckAccess(
	logger mrlog.LiteLogger,
	section mraccess.Section,
	memberProvider mraccess.MemberProvider,
	memberGroups mraccess.RightsGetter,
	rightsAvailability mraccess.RightsAvailability,
) PrepareHandlerFunc {
	if section.Privilege() != mrserver.PrivilegePublic && !rightsAvailability.HasPrivilege(section.Privilege()) {
		logger.Warn(
			fmt.Sprintf(
				"Privilege '%s' is not registered for section '%s', perhaps, it is not registered in the config or is not associated with any role",
				section.Privilege(), section.Name(),
			),
		)
	}

	if memberProvider == nil {
		logger.Error("MemberProvider is not set for section", "section", section.Name(), "error", mr.ErrInternalNilPointer.New())

		return func(handler mrserver.HttpHandler) mrserver.HttpHandler {
			handler.Func = func(_ http.ResponseWriter, _ *http.Request) error {
				return mr.ErrHttpClientUnauthorized.New()
			}

			return handler
		}
	}

	return func(handler mrserver.HttpHandler) mrserver.HttpHandler {
		handler.URL = section.BuildPath(handler.URL)

		if section.Privilege() == mrserver.PrivilegePublic && handler.Permission == mrserver.PermissionAnyUser {
			return handler
		}

		if section.Privilege() == mrserver.PrivilegePublic && handler.Permission == mrserver.PermissionGuestOnly {
			middlewareFunc := mrserver.MiddlewareHandlerCheckAccessToken(
				logger.ContextLogger(),
				handler.URL,
			)

			handler.Func = middlewareFunc(handler.Func)

			return handler
		}

		if section.Privilege() != mrserver.PrivilegePublic &&
			(handler.Permission == mrserver.PermissionAnyUser || handler.Permission == mrserver.PermissionGuestOnly) {
			logger.Warn(
				"This permission cannot be present in the private section",
				"error", fmt.Errorf("permission '%s' for handler '%s %s'", handler.Permission, handler.Method, handler.URL),
			)

			handler.Func = func(_ http.ResponseWriter, _ *http.Request) error {
				return mr.ErrHttpAccessForbidden.New()
			}

			return handler
		}

		if !rightsAvailability.HasPermission(handler.Permission) {
			logger.Warn(
				"Permission is not registered, perhaps, it is not registered in the config or is not associated with any role",
				"error", fmt.Errorf("permission '%s' for handler '%s %s'", handler.Permission, handler.Method, handler.URL),
			)

			handler.Func = func(_ http.ResponseWriter, _ *http.Request) error {
				return mr.ErrHttpAccessForbidden.New()
			}

			return handler
		}

		middlewareFunc := mrserver.MiddlewareHandlerCheckAccess(
			logger.ContextLogger(),
			handler.URL,
			section.Privilege(),
			handler.Permission,
			memberProvider,
			memberGroups,
		)

		handler.Func = middlewareFunc(handler.Func)

		return handler
	}
}
