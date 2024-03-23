package mrjulienrouter

import (
	"fmt"
	"net/http"
	"reflect"
	"runtime"

	"github.com/julienschmidt/httprouter"
	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrserver"
)

// go get -u github.com/julienschmidt/httprouter

type (
	RouterAdapter struct {
		router             *httprouter.Router
		generalHandler     http.Handler
		handlerAdapterFunc func(next mrserver.HttpHandlerFunc) http.HandlerFunc
		logger             mrlog.Logger
	}
)

// Make sure the RouterAdapter conforms with the mrserver.HttpRouter interface
var _ mrserver.HttpRouter = (*RouterAdapter)(nil)

func New(
	logger mrlog.Logger,
	handlerAdapterFunc func(next mrserver.HttpHandlerFunc) http.HandlerFunc,
	handlerNotFoundFunc http.HandlerFunc,
	handlerMethodNotAllowedFunc http.HandlerFunc,
) *RouterAdapter {
	router := httprouter.New()

	if handlerNotFoundFunc != nil {
		router.NotFound = handlerNotFoundFunc
	}

	if handlerMethodNotAllowedFunc != nil {
		router.MethodNotAllowed = handlerMethodNotAllowedFunc
	}

	return &RouterAdapter{
		router:             router,
		generalHandler:     router,
		handlerAdapterFunc: handlerAdapterFunc,
		logger:             logger,
	}
}

func (rt *RouterAdapter) RegisterMiddleware(handlers ...func(next http.Handler) http.Handler) {
	// recursion call: handler1(handler2(handler3(router())))
	for i := len(handlers) - 1; i >= 0; i-- {
		rt.generalHandler = handlers[i](rt.generalHandler)

		rt.logger.Debug().MsgFunc(
			func() string {
				return fmt.Sprintf(
					"Registered Middleware %s",
					runtime.FuncForPC(reflect.ValueOf(rt.generalHandler).Pointer()).Name(),
				)
			},
		)
	}
}

func (rt *RouterAdapter) Register(controllers ...mrserver.HttpController) {
	for i := range controllers {
		for _, handler := range controllers[i].Handlers() {
			rt.HandlerFunc(handler.Method, handler.URL, rt.handlerAdapterFunc(handler.Func))
		}
	}
}

func (rt *RouterAdapter) HandlerFunc(method, path string, handler http.HandlerFunc) {
	rt.logger.Debug().Msgf("- registered: %s %s", method, path)
	rt.router.Handler(method, path, handler)
}

func (rt *RouterAdapter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rt.generalHandler.ServeHTTP(w, r)
}
