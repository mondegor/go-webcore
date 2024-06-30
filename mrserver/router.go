package mrserver

import (
	"net/http"
)

const (
	VarRestOfURL = "{{restOfUrl}}" // VarRestOfURL - переменная остатка пути
)

type (
	// HttpRouter - роутинг запросов с их регистрацией и запуском.
	HttpRouter interface {
		RegisterMiddleware(handlers ...func(next http.Handler) http.Handler)
		Register(controllers ...HttpController)
		HandlerFunc(method, path string, handler http.HandlerFunc)
		ServeHTTP(w http.ResponseWriter, r *http.Request)
	}

	// HttpController - http контроллер со списком его обработчиков.
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
