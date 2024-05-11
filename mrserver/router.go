package mrserver

import (
	"net/http"
)

const (
	VarRestOfURL = "{restOfUrl}"
)

type (
	HTTPRouter interface {
		RegisterMiddleware(handlers ...func(next http.Handler) http.Handler)
		Register(controllers ...HTTPController)
		HandlerFunc(method, path string, handler http.HandlerFunc)
		ServeHTTP(w http.ResponseWriter, r *http.Request)
	}

	HTTPController interface {
		Handlers() []HTTPHandler
	}

	HTTPHandler struct {
		Method     string
		URL        string
		Permission string
		Func       HTTPHandlerFunc
	}

	HTTPHandlerFunc func(w http.ResponseWriter, r *http.Request) error
)
