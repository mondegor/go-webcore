package mrserver

import (
	"net/http"
)

type (
	HttpRouter interface {
		RegisterMiddleware(handlers ...HttpMiddleware)
		Register(controllers ...HttpController)
		HandlerFunc(method, path string, handler http.HandlerFunc)
		HttpHandlerFunc(method, path string, handler HttpHandlerFunc)
		ServeHTTP(w http.ResponseWriter, r *http.Request)
	}

	HttpMiddleware interface {
		Middleware(next http.Handler) http.Handler
	}

	HttpMiddlewareFunc func(next http.Handler) http.Handler

	HttpController interface {
		Handlers() []HttpHandler
	}

	HttpHandler struct {
		Method     string
		URL        string
		Permission string
		Func       HttpHandlerFunc
	}

	HttpHandlerFunc        func(w http.ResponseWriter, r *http.Request) error
	HttpHandlerAdapterFunc func(next HttpHandlerFunc) http.HandlerFunc
)

func (f HttpMiddlewareFunc) Middleware(next http.Handler) http.Handler {
	return f(next)
}
