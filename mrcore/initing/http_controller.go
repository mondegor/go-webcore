package initing

import "github.com/mondegor/go-webcore/mrserver"

type (
	// httpController - предназначен для реализации интерфейса http контроллера.
	httpController []mrserver.HttpHandler
)

// Handlers - реализует интерфейс http контроллера.
func (c httpController) Handlers() []mrserver.HttpHandler {
	return c
}

// PrepareHttpController - модифицирует обработчики контроллера, применяя к каждому
// из них переданные функции обработки, и возвращает новый контроллер с этими обработчиками.
func PrepareHttpController(c mrserver.HttpController, operations ...PrepareHandlerFunc) mrserver.HttpController {
	if len(operations) == 0 {
		return c
	}

	handlers := c.Handlers()

	for j := range handlers {
		for _, apply := range operations {
			handlers[j] = apply(handlers[j])
		}
	}

	return httpController(handlers)
}
