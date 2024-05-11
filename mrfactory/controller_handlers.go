package mrfactory

import (
	"github.com/mondegor/go-webcore/mrserver"
)

type (
	ControllerHandlers      []mrserver.HTTPHandler
	PrepareHandlerFunc      func(handler *mrserver.HTTPHandler)
	ApplyEachControllerFunc func(list []mrserver.HTTPController, err error) error
)

func (c ControllerHandlers) Handlers() []mrserver.HTTPHandler {
	return c
}

func PrepareEachController(list []mrserver.HTTPController, operations ...PrepareHandlerFunc) []mrserver.HTTPController {
	for i := range list {
		list[i] = PrepareController(list[i], operations...)
	}

	return list
}

func PrepareController(c mrserver.HTTPController, operations ...PrepareHandlerFunc) mrserver.HTTPController {
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
