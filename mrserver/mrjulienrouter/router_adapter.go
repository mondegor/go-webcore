package mrjulienrouter

import (
	"net/http"
	"reflect"
	"runtime"

	"github.com/julienschmidt/httprouter"

	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrserver"
)

// go get -u github.com/julienschmidt/httprouter

const (
	varRestOfURL = "path"
)

type (
	// RouterAdapter - comment struct.
	RouterAdapter struct {
		router             *httprouter.Router
		generalHandler     http.Handler
		handlerAdapterFunc func(next mrserver.HttpHandlerFunc) http.HandlerFunc
		logger             mrlog.Logger
	}
)

// Make sure the RouterAdapter conforms with the mrserver.HttpRouter interface.
var _ mrserver.HttpRouter = (*RouterAdapter)(nil)

// New - создаёт объект RouterAdapter.
func New(
	logger mrlog.Logger,
	adapterFunc func(next mrserver.HttpHandlerFunc) http.HandlerFunc,
	notFoundFunc http.HandlerFunc,
	methodNotAllowedFunc http.HandlerFunc,
) *RouterAdapter {
	router := httprouter.New()

	if notFoundFunc != nil {
		router.NotFound = notFoundFunc
	}

	if methodNotAllowedFunc != nil {
		router.MethodNotAllowed = methodNotAllowedFunc
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
	convertedPath := ConvertURL(path)

	if path != convertedPath {
		path += " -> " + convertedPath
	}

	rt.logger.Debug().Msgf("- registered: %s %s", method, path)
	rt.router.Handler(method, convertedPath, handler)
}

// ServeHTTP - comment method.
func (rt *RouterAdapter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rt.generalHandler.ServeHTTP(w, r)
}
