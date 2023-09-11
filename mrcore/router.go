package mrcore

import "net/http"

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
        AddHandlers(router HttpRouter)
    }

    HttpHandlerFunc func(c ClientData) error
    HttpHandlerAdapterFunc func(next HttpHandlerFunc) http.HandlerFunc
)

func (f HttpMiddlewareFunc) Middleware(next http.Handler) http.Handler {
    return f(next)
}
