package mrserver

import (
    "net/http"
    "reflect"
    "runtime"

    "github.com/julienschmidt/httprouter"
    "github.com/mondegor/go-webcore/mrcore"
)

// go get -u github.com/julienschmidt/httprouter

// Make sure the routerAdapter conforms with the mrcore.HttpRouter interface
var _ mrcore.HttpRouter = (*routerAdapter)(nil)

type (
    routerAdapter struct {
        router *httprouter.Router
        generalHandler http.Handler
        handlerAdapterFunc mrcore.HttpHandlerAdapterFunc
        logger mrcore.Logger
    }
)

func NewRouter(
    logger mrcore.Logger,
    handlerAdapterFunc mrcore.HttpHandlerAdapterFunc,
) *routerAdapter {
    router := httprouter.New()

    // r.GlobalOPTIONS
    // rt.router.MethodNotAllowed
    // rt.router.NotFound

    return &routerAdapter{
        router: router,
        generalHandler: router,
        handlerAdapterFunc: handlerAdapterFunc,
        logger: logger,
    }
}

func (rt *routerAdapter) RegisterMiddleware(handlers ...mrcore.HttpMiddleware) {
    // recursion call: handler1(handler2(handler3(router())))
    for i := len(handlers) - 1; i >= 0; i-- {
        rt.generalHandler = handlers[i].Middleware(rt.generalHandler)

        rt.logger.Info(
            "Registered Middleware %s",
            runtime.FuncForPC(reflect.ValueOf(rt.generalHandler).Pointer()).Name(),
        )
    }
}

func (rt *routerAdapter) Register(controllers ...mrcore.HttpController) {
    for _, controller := range controllers {
        controller.AddHandlers(rt)
    }
}

func (rt *routerAdapter) HandlerFunc(method, path string, handler http.HandlerFunc) {
    rt.router.Handler(method, path, handler)
}

func (rt *routerAdapter) HttpHandlerFunc(method, path string, handler mrcore.HttpHandlerFunc) {
    rt.router.Handler(method, path, rt.handlerAdapterFunc(handler))
}

func (rt *routerAdapter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    rt.generalHandler.ServeHTTP(w, r)
}
