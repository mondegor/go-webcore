package mrresp

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrlib"
	"github.com/mondegor/go-webcore/mrlog"
)

// HandlerGetNotFoundAsJSON  - comment func.
func HandlerGetNotFoundAsJSON(unexpectedStatus int) http.HandlerFunc {
	return HandlerErrorResponse(
		http.StatusNotFound,
		unexpectedStatus,
		"404 Not Found",
		"The server cannot find the requested resource",
	)
}

// HandlerGetMethodNotAllowedAsJSON  - comment func.
func HandlerGetMethodNotAllowedAsJSON(unexpectedStatus int) http.HandlerFunc {
	return HandlerErrorResponse(
		http.StatusMethodNotAllowed,
		unexpectedStatus,
		"405 Method Not Allowed",
		"The server knows the request method, but the target resource doesn't support this method",
	)
}

// HandlerGetFatalErrorAsJSON  - comment func.
func HandlerGetFatalErrorAsJSON(unexpectedStatus int) http.HandlerFunc {
	return HandlerErrorResponse(
		http.StatusInternalServerError,
		unexpectedStatus,
		"Internal server error",
		"The server encountered an unexpected condition that prevented it from fulfilling the request",
	)
}

// HandlerErrorResponse  - comment func.
func HandlerErrorResponse(status, unexpectedStatus int, title, details string) http.HandlerFunc {
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
			status = unexpectedStatus
			bytes = nil

			mrlog.Ctx(r.Context()).Error().Err(mrcore.ErrHttpResponseParseData.Wrap(err)).Msg("marshal failed")
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		mrlib.Write(r.Context(), w, bytes)
	}
}
