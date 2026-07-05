package mrresp

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/mondegor/go-sysmess/errors"
	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-sysmess/util/xio"
)

// HandlerGetNotFoundAsJSON - создаёт обработчик для ответов 404 Not Found.
// Формирует ответ в формате RFC 9457 (Problem Details).
func HandlerGetNotFoundAsJSON(logger mrlog.Logger) http.HandlerFunc {
	return HandlerErrorResponse(
		logger,
		http.StatusNotFound,
		"404 Not Found",
		"The server cannot find the requested resource",
	)
}

// HandlerGetMethodNotAllowedAsJSON - создаёт обработчик для ответов 405 Method Not Allowed.
// Формирует ответ в формате RFC 9457 (Problem Details).
func HandlerGetMethodNotAllowedAsJSON(logger mrlog.Logger) http.HandlerFunc {
	return HandlerErrorResponse(
		logger,
		http.StatusMethodNotAllowed,
		"405 Method Not Allowed",
		"The server knows the request method, but the target resource doesn't support this method",
	)
}

// HandlerGetFatalErrorAsJSON - создаёт обработчик для ответов 500 Internal Server Error.
// Формирует ответ в формате RFC 9457 (Problem Details).
func HandlerGetFatalErrorAsJSON(logger mrlog.Logger) http.HandlerFunc {
	return HandlerErrorResponse(
		logger,
		http.StatusInternalServerError,
		"Internal server error",
		"The server encountered an unexpected condition that prevented it from fulfilling the request",
	)
}

// HandlerErrorResponse - создаёт универсальный обработчик для ответов с ошибкой.
// Формирует ответ в формате RFC 9457 (Problem Details for HTTP APIs).
// Instance и Time заполняются автоматически из запроса.
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
