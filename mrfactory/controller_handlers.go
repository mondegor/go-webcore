package mrfactory

import (
	"github.com/mondegor/go-webcore/mrserver"
)

type (
	ControllerHandlers      []mrserver.HttpHandler
	PrepareHandlerFunc      func(handler *mrserver.HttpHandler)
	ApplyEachControllerFunc func(list []mrserver.HttpController, err error) error
)

func (c ControllerHandlers) Handlers() []mrserver.HttpHandler {
	return c
}

func PrepareEachController(list []mrserver.HttpController, operations ...PrepareHandlerFunc) []mrserver.HttpController {
	for i := range list {
		list[i] = PrepareController(list[i], operations...)
	}

	return list
}

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
