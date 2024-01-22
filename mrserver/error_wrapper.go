package mrserver

import (
	"net/http"

	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-webcore/mrcore"
)

// DefaultHttpErrorOverrideFunc - only for: 401, 403, 404, 418, 5XX
func DefaultHttpErrorOverrideFunc(err *mrerr.AppError) (int, *mrerr.AppError) {
	status := http.StatusInternalServerError

	if mrcore.FactoryErrServiceEntityNotFound.Is(err) ||
		mrcore.FactoryErrHttpResourceNotFound.Is(err) {
		status = http.StatusNotFound
	} else if mrcore.FactoryErrHttpClientUnauthorized.Is(err) {
		status = http.StatusUnauthorized
	} else if mrcore.FactoryErrHttpAccessForbidden.Is(err) {
		status = http.StatusForbidden
	} else if err.ID() == mrcore.FactoryErrInternal.ErrorID() {
		// если ошибка явно не обработана разработчиком, то вместо 500 отображается 418
		status = http.StatusTeapot
	}

	return status, err
}
