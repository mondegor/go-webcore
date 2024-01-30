package mrserver

import (
	"net/http"

	"github.com/mondegor/go-webcore/mrcore"
)

func NewMiddlewareHttpHandlerAdapter(s ErrorResponseSender) (HttpHandlerAdapterFunc, error) {
	if s == nil {
		return nil, mrcore.FactoryErrInternalNilPointer.New()
	}

	return func(next HttpHandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if err := next(w, r); err != nil {
				s.SendError(w, r, err)
			}
		}
	}, nil
}
