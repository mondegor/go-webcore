package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"runtime/debug"

	"github.com/mondegor/go-sysmess/mrerr/mr"
	"github.com/mondegor/go-sysmess/mrlog"
)

// RecoverHandler - промежуточный обработчик для перехвата panic.
func RecoverHandler(logger mrlog.Logger, isDebug bool, fatalFunc http.HandlerFunc) func(next http.Handler) http.Handler {
	if fatalFunc == nil {
		fatalFunc = func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func(ctx context.Context) {
				if rvr := recover(); rvr != nil {
					if rvr == http.ErrAbortHandler { //nolint:errorlint
						// we don't recover http.ErrAbortHandler so the response
						// to the client is aborted, this should not be logged
						panic(rvr)
					}

					errorMessage := fmt.Sprintf("proto='%s', method='%s', url='%s'", r.Proto, r.Method, r.URL)

					if isDebug {
						os.Stderr.Write([]byte(errorMessage))
						os.Stderr.Write([]byte(fmt.Sprintf("; panic: %v\n", rvr)))
						os.Stderr.Write(debug.Stack())
					} else {
						logger.Error(
							ctx,
							"RecoverHandler",
							"error",
							mr.ErrInternalCaughtPanic.New(
								errorMessage,
								rvr,
								string(debug.Stack()),
							),
						)
					}

					fatalFunc(w, r)
				}
			}(r.Context())

			next.ServeHTTP(w, r)
		})
	}
}
