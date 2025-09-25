package mrchi

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/mondegor/go-webcore/mrserver"
)

// URLPathParam - comment func.
func URLPathParam(r *http.Request, name string) string {
	if ctx, ok := r.Context().Value(chi.RouteCtxKey).(*chi.Context); ok {
		if name == mrserver.VarRestOfURL {
			name = "*"
		}

		return ctx.URLParam(name)
	}

	// :TODO: переделать в объект
	// mrlog.Ctx(r.Context()).Warn().Msg("chi.RouteCtxKey is not found in context")

	return ""
}
