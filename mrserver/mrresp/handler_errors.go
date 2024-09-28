package mrresp

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrlib"
	"github.com/mondegor/go-webcore/mrlog"
)

// HandlerGetNotFoundAsJSON - возвращает обработчик для формирования 404 ошибки.
func HandlerGetNotFoundAsJSON() http.HandlerFunc {
	return HandlerErrorResponse(
		http.StatusNotFound,
		"404 Not Found",
		"The server cannot find the requested resource",
	)
}

// HandlerGetMethodNotAllowedAsJSON - возвращает обработчик для формирования 405 ошибки.
func HandlerGetMethodNotAllowedAsJSON() http.HandlerFunc {
	return HandlerErrorResponse(
		http.StatusMethodNotAllowed,
		"405 Method Not Allowed",
		"The server knows the request method, but the target resource doesn't support this method",
	)
}

// HandlerGetFatalErrorAsJSON - возвращает обработчик для формирования 500 ошибки.
func HandlerGetFatalErrorAsJSON() http.HandlerFunc {
	return HandlerErrorResponse(
		http.StatusInternalServerError,
		"Internal server error",
		"The server encountered an unexpected condition that prevented it from fulfilling the request",
	)
}

// HandlerErrorResponse - возвращает обработчик для формирования ошибки согласно RFC 7807 (Problem Details for HTTP APIs).
func HandlerErrorResponse(status int, title, details string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bytes, err := json.Marshal(
			ErrorDetailsResponse{
				Title:   title,
				Details: details,
				Request: r.URL.Path,
				Time:    time.Now().UTC().Format(time.RFC3339),
			},
		)
		if err != nil {
			status = http.StatusUnprocessableEntity
			bytes = nil

			mrlog.Ctx(r.Context()).Error().Err(mrcore.ErrHttpResponseParseData.Wrap(err)).Msg("marshal failed")
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		mrlib.Write(r.Context(), w, bytes)
	}
}
