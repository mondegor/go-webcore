package mrserver

import (
	"net/http"

	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-webcore/mrcore"
)

// DefaultHTTPErrorOverrideFunc - only for: 401, 403, 404, 418, 5XX
func DefaultHTTPErrorOverrideFunc(err *mrerr.AppError) (int, *mrerr.AppError) {
	status := http.StatusInternalServerError

	if mrcore.FactoryErrUseCaseEntityNotFound.Is(err) ||
		mrcore.FactoryErrHTTPResourceNotFound.Is(err) {
		status = http.StatusNotFound
	} else if mrcore.FactoryErrHTTPClientUnauthorized.Is(err) {
		status = http.StatusUnauthorized
	} else if mrcore.FactoryErrHTTPAccessForbidden.Is(err) {
		status = http.StatusForbidden
	} else if err.Code() == mrcore.FactoryErrInternal.Code() {
		// если ошибка явно не обработана разработчиком, то вместо 500 отображается 418
		status = http.StatusTeapot
	}

	return status, err
}
