package mrfactory

import (
	"github.com/mondegor/go-webcore/mrserver"
)

type (
	// ControllerHandlers - обработчики контроллера.
	ControllerHandlers []mrserver.HttpHandler

	// PrepareHandlerFunc - функция для подготовки обработчика (URL, разрешения).
	PrepareHandlerFunc func(handler *mrserver.HttpHandler)

	// ApplyEachControllerFunc - вспомогательная функция для подготовки списка контроллеров.
	ApplyEachControllerFunc func(list []mrserver.HttpController, err error) error
)

// Handlers - comment method.
func (c ControllerHandlers) Handlers() []mrserver.HttpHandler {
	return c
}

// PrepareEachController - comment func.
func PrepareEachController(list []mrserver.HttpController, operations ...PrepareHandlerFunc) []mrserver.HttpController {
	for i := range list {
		list[i] = PrepareController(list[i], operations...)
	}

	return list
}

// PrepareController - comment func.
func PrepareController(c mrserver.HttpController, operations ...PrepareHandlerFunc) mrserver.HttpController {
	if len(operations) == 0 {
		return c
	}

	handlers := c.Handlers()

	for j := range handlers {
		for _, fn := range operations {
			fn(&handlers[j])
		}
	}

	return ControllerHandlers(handlers)
}
