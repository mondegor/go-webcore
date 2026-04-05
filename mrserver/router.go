package mrserver

import (
	"net/http"
)

const (
	// VarRestOfURL - переменная остатка пути.
	VarRestOfURL = "{{restOfUrl}}"
)

type (
	// HttpRouter - обеспечивает регистрацию middleware,
	// контроллеров и маршрутизацию HTTP-запросов.
	HttpRouter interface {
		RegisterMiddleware(handlers ...func(next http.Handler) http.Handler)
		Register(controllers ...HttpController)
		HandlerFunc(method, path string, handler http.HandlerFunc)
		ServeHTTP(w http.ResponseWriter, r *http.Request)
	}

	// HttpController - определяет контроллер с набором HTTP-обработчиков.
	HttpController interface {
		Handlers() []HttpHandler
	}

	// HttpHandler - http обработчик, к которому привязаны метод, URL и разрешение
	// контролирующее запуск этого обработчика.
	HttpHandler struct {
		Method     string
		URL        string
		Permission string
		Func       HttpHandlerFunc
	}

	// HttpHandlerFunc - изменённый дизайн стандартного HTTP обработчика
	// с возможностью возврата ошибки вместо её обработки в самом обработчике.
	HttpHandlerFunc func(w http.ResponseWriter, r *http.Request) error
)
