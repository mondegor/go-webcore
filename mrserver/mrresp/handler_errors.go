package mrresp

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/mondegor/go-sysmess/errors"
	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-sysmess/util/xio"
)

// HandlerGetNotFoundAsJSON - возвращает обработчик для формирования 404 ошибки.
func HandlerGetNotFoundAsJSON(logger mrlog.Logger) http.HandlerFunc {
	return HandlerErrorResponse(
		logger,
		http.StatusNotFound,
		"404 Not Found",
		"The server cannot find the requested resource",
	)
}

// HandlerGetMethodNotAllowedAsJSON - возвращает обработчик для формирования 405 ошибки.
func HandlerGetMethodNotAllowedAsJSON(logger mrlog.Logger) http.HandlerFunc {
	return HandlerErrorResponse(
		logger,
		http.StatusMethodNotAllowed,
		"405 Method Not Allowed",
		"The server knows the request method, but the target resource doesn't support this method",
	)
}

// HandlerGetFatalErrorAsJSON - возвращает обработчик для формирования 500 ошибки.
func HandlerGetFatalErrorAsJSON(logger mrlog.Logger) http.HandlerFunc {
	return HandlerErrorResponse(
		logger,
		http.StatusInternalServerError,
		"Internal server error",
		"The server encountered an unexpected condition that prevented it from fulfilling the request",
	)
}

// HandlerErrorResponse - возвращает обработчик для формирования ошибки согласно RFC 9457 (Problem Details for HTTP APIs).
func HandlerErrorResponse(logger mrlog.Logger, status int, title, detail string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bytes, err := json.Marshal(
			ErrorDetailsResponse{
				Title:    title,
				Status:   status,
				Detail:   detail,
				Instance: r.Method + " " + r.URL.Path, // TODO: добавить helper xhttp.RequestInstance(r)
				Time:     time.Now().UTC().Format(time.RFC3339),
			},
		)
		if err != nil {
			status = http.StatusUnprocessableEntity
			bytes = nil

			logger.Error(r.Context(), "marshal failed", "error", errors.ErrInternalHttpResponseParseData.Wrap(err))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		xio.Write(r.Context(), logger, w, bytes)
	}
}
