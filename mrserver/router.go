package mrserver

import (
	"net/http"
)

type (
	HttpRouter interface {
		RegisterMiddleware(handlers ...func(next http.Handler) http.Handler)
		Register(controllers ...HttpController)
		HandlerFunc(method, path string, handler http.HandlerFunc)
		ServeHTTP(w http.ResponseWriter, r *http.Request)
	}

	HttpController interface {
		Handlers() []HttpHandler
	}

	HttpHandler struct {
		Method     string
		URL        string
		Permission string
		Func       HttpHandlerFunc
	}

	HttpHandlerFunc func(w http.ResponseWriter, r *http.Request) error
)
