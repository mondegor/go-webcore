package initing

import "github.com/mondegor/go-webcore/mrserver"

type (
	httpController []mrserver.HttpHandler
)

// Handlers - возвращает список HTTP-обработчиков контроллера.
// Реализует интерфейс mrserver.HttpController.
func (c httpController) Handlers() []mrserver.HttpHandler {
	return c
}

// PrepareHttpController - применяет цепочку преобразований к обработчикам контроллера.
// Параметры:
//   - c - исходный контроллер с обработчиками;
//   - operations - функции преобразования, применяемые последовательно к каждому обработчику;
//
// Возвращает новый контроллер с модифицированными обработчиками.
// Если operations пуст, возвращает исходный контроллер без изменений.
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
