package mrserver

import (
	"errors"
	"net/http"

	"github.com/mondegor/go-webcore/mrcore"
)

type (
	// HttpErrorStatusGetter - only for: 401, 403, 404, 418, 5XX.
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

// ErrorStatus - comment method.
func (g *HttpErrorStatusGetter) ErrorStatus(err error) int {
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

	status := http.StatusInternalServerError

	// если ошибка явно не обработана разработчиком (ни чем не обёрнута),
	// то вместо 500 отображается указанный g.unexpectedStatus
	if g.unexpectedStatus != status && mrcore.IsUnexpectedError(err) {
		return g.unexpectedStatus
	}

	return status
}
