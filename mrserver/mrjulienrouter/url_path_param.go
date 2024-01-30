package mrjulienrouter

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/mondegor/go-webcore/mrlog"
)

func PathParam(r *http.Request, name string) string {
	if params, ok := r.Context().Value(httprouter.ParamsKey).(httprouter.Params); ok {
		return params.ByName(name)
	}

	mrlog.Ctx(r.Context()).Warn().Msg("httprouter.ParamsKey is not found in context")

	return ""
}
