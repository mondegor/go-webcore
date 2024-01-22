package mrfactory

import (
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrserver"
)

func WithPermission(list []mrserver.HttpController, permission string) []mrserver.HttpController {
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
	list []mrserver.HttpController,
	section mrcore.AppSection,
	access mrcore.AccessControl,
) []mrserver.HttpController {
	for i := range list {
		handlers := list[i].Handlers()

		for j := range handlers {
			if !access.HasPermission(handlers[j].Permission) {
				mrcore.LogWarning(
					"permission '%s' is not registered for handler '%s %s', perhaps, it is not registered in the config or is not associated with any role",
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
