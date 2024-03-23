package mrfactory

import (
	"context"

	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrperms"
	"github.com/mondegor/go-webcore/mrserver"
)

func WithPermission(permission string) PrepareHandlerFunc {
	return func(handler *mrserver.HttpHandler) {
		if handler.Permission == "" {
			handler.Permission = permission
		}
	}
}

func WithMiddlewareCheckAccess(ctx context.Context, section mrperms.AppSection, access mrperms.AccessControl) PrepareHandlerFunc {
	return func(handler *mrserver.HttpHandler) {
		if !access.HasPermission(handler.Permission) {
			mrlog.Ctx(ctx).Warn().Caller(1).Msgf(
				"Permission '%s' is not registered for handler '%s %s', perhaps, it is not registered in the config or is not associated with any role",
				handler.Permission,
				handler.Method,
				handler.URL,
			)
		}

		fn := mrserver.MiddlewareHandlerCheckAccess(
			handler.URL,
			access,
			section.Privilege(),
			handler.Permission,
		)

		handler.URL = section.Path(handler.URL)
		handler.Func = fn(handler.Func)
	}
}
