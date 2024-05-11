package mrjulienrouter

import (
	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"runtime"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrserver"
)

// go get -u github.com/julienschmidt/httprouter

const (
	varRestOfURL = "path"
)

type (
	RouterAdapter struct {
		router             *httprouter.Router
		generalHandler     http.Handler
		handlerAdapterFunc func(next mrserver.HTTPHandlerFunc) http.HandlerFunc
		logger             mrlog.Logger
	}
)

var (
	// Make sure the RouterAdapter conforms with the mrserver.HTTPRouter interface
	_ mrserver.HTTPRouter = (*RouterAdapter)(nil)

	regexpURLVars = regexp.MustCompile(`{([a-zA-Z][a-zA-Z0-9_]*)}`)
)

func New(
	logger mrlog.Logger,
	adapterFunc func(next mrserver.HTTPHandlerFunc) http.HandlerFunc,
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

func (rt *RouterAdapter) Register(controllers ...mrserver.HTTPController) {
	for i := range controllers {
		for _, handler := range controllers[i].Handlers() {
			rt.HandlerFunc(handler.Method, handler.URL, rt.handlerAdapterFunc(handler.Func))
		}
	}
}

func (rt *RouterAdapter) HandlerFunc(method, path string, handler http.HandlerFunc) {
	convertedPath := rt.convertURL(path)

	if path != convertedPath {
		path += " -> " + convertedPath
	}

	rt.logger.Debug().Msgf("- registered: %s %s", method, path)
	rt.router.Handler(method, convertedPath, handler)
}

func (rt *RouterAdapter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rt.generalHandler.ServeHTTP(w, r)
}

func (rt *RouterAdapter) convertURL(url string) string {
	url = strings.Replace(url, mrserver.VarRestOfURL, "*"+varRestOfURL, 1)

	for _, m := range regexpURLVars.FindAllStringSubmatch(url, -1) {
		url = strings.Replace(url, m[0], ":"+m[1], 1)
	}

	return url
}
