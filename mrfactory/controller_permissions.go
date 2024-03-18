package mrfactory

import (
	"context"

	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrperms"
	"github.com/mondegor/go-webcore/mrserver"
)

func WithPermission(ctx context.Context, list []mrserver.HttpController, permission string) []mrserver.HttpController {
	for i := range list {
		handlers := list[i].Handlers()

		for j := range handlers {
			if handlers[j].Permission == "" {
				handlers[j].Permission = permission
			}
		}

		list[i] = ControllerHandlers(handlers)
	}

	return list
}

func WithMiddlewareCheckAccess(
	ctx context.Context,
	list []mrserver.HttpController,
	section mrperms.AppSection,
	access mrperms.AccessControl,
) []mrserver.HttpController {
	for i := range list {
		handlers := list[i].Handlers()

		for j := range handlers {
			if !access.HasPermission(handlers[j].Permission) {
				mrlog.Ctx(ctx).Warn().Caller(1).Msgf(
					"Permission '%s' is not registered for handler '%s %s', perhaps, it is not registered in the config or is not associated with any role",
					handlers[j].Permission,
					handlers[j].Method,
					handlers[j].URL,
				)
			}

			handlers[j].URL = section.Path(handlers[j].URL)
			handlers[j].Func = mrserver.MiddlewareCheckAccess(
				section,
				access,
				handlers[j].Permission,
				handlers[j].Func,
			)
		}

		list[i] = ControllerHandlers(handlers)
	}

	return list
}
