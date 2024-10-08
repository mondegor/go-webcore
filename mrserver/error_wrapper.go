package mrserver

import (
	"errors"
	"net/http"

	"github.com/mondegor/go-webcore/mrcore"
)

type (
	// HttpErrorStatusGetter - only for: 401, 403, 404, 418, 422, 5XX.
	HttpErrorStatusGetter struct {
		unexpectedStatus int
	}
)

// NewHttpErrorStatusGetter - создаёт объект HttpErrorStatusGetter.
func NewHttpErrorStatusGetter(unexpectedStatus int) *HttpErrorStatusGetter {
	return &HttpErrorStatusGetter{
		unexpectedStatus: unexpectedStatus,
	}
}

// ErrorStatus - возвращает http код ответа на основе типа ошибки и самой ошибки.
func (g *HttpErrorStatusGetter) ErrorStatus(errorType mrcore.AnalyzedErrorType, err error) int {
	if errors.Is(err, mrcore.ErrUseCaseEntityNotFound) ||
		errors.Is(err, mrcore.ErrHttpResourceNotFound) {
		return http.StatusNotFound
	}

	if errors.Is(err, mrcore.ErrHttpClientUnauthorized) {
		return http.StatusUnauthorized
	}

	if errors.Is(err, mrcore.ErrHttpAccessForbidden) {
		return http.StatusForbidden
	}

	if errors.Is(err, mrcore.ErrHttpRequestParseData) {
		return http.StatusUnprocessableEntity
	}

	if errorType == mrcore.AnalyzedErrorTypeUser || errorType == mrcore.AnalyzedErrorTypeProtoUser {
		return http.StatusBadRequest
	}

	if errorType == mrcore.AnalyzedErrorTypeSystem || errorType == mrcore.AnalyzedErrorTypeProtoSystem {
		return http.StatusServiceUnavailable
	}

	// если ошибка явно необработанна разработчиком (не обёрнута в ProtoAppError),
	// то вместо 500 статуса отображается указанный g.unexpectedStatus
	if errorType == mrcore.AnalyzedErrorTypeUndefined {
		return g.unexpectedStatus
	}

	return http.StatusInternalServerError
}
