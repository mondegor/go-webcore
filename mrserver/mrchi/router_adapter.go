package mrchi

import (
	"net/http"
	"reflect"
	"runtime"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrserver"
)

// go get -u github.com/go-chi/chi/v5

type (
	// RouterAdapter - comment struct.
	RouterAdapter struct {
		router             *chi.Mux
		generalHandler     http.Handler
		handlerAdapterFunc func(next mrserver.HttpHandlerFunc) http.HandlerFunc
		logger             mrlog.Logger
	}
)

// New - создаёт объект RouterAdapter.
func New(
	logger mrlog.Logger,
	adapterFunc func(next mrserver.HttpHandlerFunc) http.HandlerFunc,
	notFoundFunc http.HandlerFunc,
	methodNotAllowedFunc http.HandlerFunc,
) *RouterAdapter {
	router := chi.NewRouter()

	if notFoundFunc != nil {
		router.NotFound(notFoundFunc)
	}

	if methodNotAllowedFunc != nil {
		router.MethodNotAllowed(methodNotAllowedFunc)
	}

	return &RouterAdapter{
		router:             router,
		generalHandler:     router,
		handlerAdapterFunc: adapterFunc,
		logger:             logger,
	}
}

// RegisterMiddleware - comment method.
func (rt *RouterAdapter) RegisterMiddleware(handlers ...func(next http.Handler) http.Handler) {
	// recursion call: handler1(handler2(handler3(router())))
	for i := len(handlers) - 1; i >= 0; i-- {
		rt.generalHandler = handlers[i](rt.generalHandler)

		rt.logger.Debug().MsgFunc(
			func() string {
				return "Registered Middleware " +
					runtime.FuncForPC(reflect.ValueOf(rt.generalHandler).Pointer()).Name()
			},
		)
	}
}

// Register - comment method.
func (rt *RouterAdapter) Register(controllers ...mrserver.HttpController) {
	for i := range controllers {
		for _, handler := range controllers[i].Handlers() {
			rt.HandlerFunc(handler.Method, handler.URL, rt.handlerAdapterFunc(handler.Func))
		}
	}
}

// HandlerFunc - comment method.
func (rt *RouterAdapter) HandlerFunc(method, path string, handler http.HandlerFunc) {
	convertedPath := rt.convertURL(path)

	if path != convertedPath {
		path += " -> " + convertedPath
	}

	rt.logger.Debug().Msgf("- registered: %s %s", method, path)
	rt.router.Method(method, convertedPath, handler)
}

// ServeHTTP - comment method.
func (rt *RouterAdapter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rt.generalHandler.ServeHTTP(w, r)
}

func (rt *RouterAdapter) convertURL(url string) string {
	return strings.Replace(url, mrserver.VarRestOfURL, "*", 1)
}
