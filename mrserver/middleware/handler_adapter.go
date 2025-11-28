package middleware

import (
	"net/http"

	"github.com/mondegor/go-sysmess/mrerr/mr"

	"github.com/mondegor/go-webcore/mrserver"
)

// HandlerAdapter - переходник с HttpHandlerFunc на http.HandlerFunc.
func HandlerAdapter(errSender mrserver.ErrorResponseSender) func(next mrserver.HttpHandlerFunc) http.HandlerFunc {
	return func(next mrserver.HttpHandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if err := next(w, r); err != nil {
				if mr.ErrUseCaseEntityNotFound.Is(err) { // подменяются только необёрнутые ошибки этого типа
					err = mr.ErrHttpResourceNotFound.New()
				}

				errSender.SendError(w, r, err)
			}
		}
	}
}
