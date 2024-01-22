package mrfactory

import (
	"github.com/mondegor/go-webcore/mrserver"
)

type (
	ControllerHandlers []mrserver.HttpHandler
)

func (c ControllerHandlers) Handlers() []mrserver.HttpHandler {
	return c
}
