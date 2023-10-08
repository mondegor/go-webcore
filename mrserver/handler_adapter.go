package mrserver

import (
    "net/http"

    "github.com/mondegor/go-webcore/mrcore"
    "github.com/mondegor/go-webcore/mrctx"
)

func HandlerAdapter(validator mrcore.Validator) mrcore.HttpHandlerAdapterFunc {
    return func(next mrcore.HttpHandlerFunc) http.HandlerFunc {
        return func(w http.ResponseWriter, r *http.Request) {
            logger := mrctx.Logger(r.Context())
            logger.Debug("Exec HandlerAdapter")

            c := clientContext{
                request: r,
                responseWriter: w,
                validator: validator,
            }

            if err := next(&c); err != nil {
                c.sendErrorResponse(err)
            }
        }
    }
}
