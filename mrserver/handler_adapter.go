package mrserver

import (
	"net/http"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrctx"
)

func HandlerAdapter(wefn mrcore.ClientWrapErrorFunc) mrcore.HttpHandlerAdapterFunc {
	return func(next mrcore.HttpHandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ct, err := mrctx.GetClientTools(r.Context())

			if err != nil {
				mrcore.LogErr(err)
				w.WriteHeader(http.StatusTeapot)
				return
			}

			ct.Logger.Debug("Exec HandlerAdapter")

			c := newClientData(r, w, wefn, ct)

			if err = next(c); err != nil {
				c.sendErrorResponse(err)
			}
		}
	}
}
