package mrchi

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/mondegor/go-webcore/mrserver"
)

// URLPathParam - извлекает параметр из URL-пути, используя маршрутизацию chi.
//
// Поддерживает:
//   - Именованные параметры chi (например: {id}, {name});
//   - Специальный параметр mrserver.VarRestOfURL (оставшаяся часть пути);
//
// Возвращает значение параметра или пустую строку, если параметр не найден.
func URLPathParam(r *http.Request, name string) string {
	if chiCtx, ok := r.Context().Value(chi.RouteCtxKey).(*chi.Context); ok {
		if name == mrserver.VarRestOfURL {
			name = "*"
		}

		return chiCtx.URLParam(name)
	}

	// :TODO: переделать в объект
	// mrlog.Ctx(r.Context()).Warn().Msg("chi.RouteCtxKey is not found in context")

	return ""
}
