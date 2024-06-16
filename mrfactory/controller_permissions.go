package mrfactory

import (
	"context"
	"fmt"

	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrperms"
	"github.com/mondegor/go-webcore/mrserver"
)

type (
	// accessControl - проверяет наличия указанных привилегий и разрешений.
	accessControl interface {
		HasPrivilege(name string) bool
		HasPermission(name string) bool
		mrperms.AccessRightsFactory
	}
)

// NewAppSection - создаёт объект mrperms.AppSection с указанными настройками.
func NewAppSection(ctx context.Context, opts mrperms.AppSectionOptions, access accessControl) *mrperms.AppSection {
	logger := mrlog.Ctx(ctx)
	logger.Info().Msgf("Init section %s with root path '%s' and privilege '%s'", opts.Caption, opts.BasePath, opts.Privilege)
	logger.Debug().Msgf("secret=%s, audience: %s", opts.AuthSecret, opts.AuthAudience)

	if !access.HasPrivilege(opts.Privilege) {
		logger.Warn().Msgf(
			"Privilege '%s' is not registered for section '%s', perhaps, it is not registered in the config or is not associated with any role",
			opts.Privilege,
			opts.Caption,
		)
	}

	return mrperms.NewAppSection(opts)
}

// WithPermission - comment func.
func WithPermission(permission string) PrepareHandlerFunc {
	return func(handler *mrserver.HttpHandler) {
		if handler.Permission == "" {
			handler.Permission = permission
		}
	}
}

// WithMiddlewareCheckAccess - comment func.
func WithMiddlewareCheckAccess(ctx context.Context, section *mrperms.AppSection, access accessControl) PrepareHandlerFunc {
	return func(handler *mrserver.HttpHandler) {
		if !access.HasPermission(handler.Permission) {
			mrlog.Ctx(ctx).Warn().Err(
				fmt.Errorf("permission '%s' for handler '%s %s'", handler.Permission, handler.Method, handler.URL),
			).Msg("Permission is not registered, perhaps, it is not registered in the config or is not associated with any role")
		}

		fn := mrserver.MiddlewareHandlerCheckAccess(
			handler.URL,
			access,
			section.Privilege(),
			handler.Permission,
		)

		handler.URL = section.BuildPath(handler.URL)
		handler.Func = fn(handler.Func)
	}
}
