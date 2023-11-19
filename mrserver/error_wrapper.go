package mrserver

import (
	"net/http"

	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-webcore/mrcore"
)

func DefaultErrorWrapperFunc(err *mrerr.AppError) (int, *mrerr.AppError) {
	status := http.StatusInternalServerError

	if mrcore.FactoryErrServiceEntityNotFound.Is(err) {
		status = http.StatusNotFound
		err = mrcore.FactoryErrHttpResourceNotFound.Wrap(err)
	} else if mrcore.FactoryErrHttpClientUnauthorized.Is(err) {
		status = http.StatusUnauthorized
	} else if mrcore.FactoryErrHttpAccessForbidden.Is(err) {
		status = http.StatusForbidden
	} else if mrcore.FactoryErrServiceEmptyInputData.Is(err) ||
		mrcore.FactoryErrServiceIncorrectInputData.Is(err) {
		err = mrcore.FactoryErrHttpRequestParseData.Wrap(err)
	} else if mrcore.FactoryErrServiceTemporarilyUnavailable.Is(err) {
		err = mrcore.FactoryErrHttpResponseSystemTemporarilyUnableToProcess.Wrap(err)
	} else if err.ID() == mrerr.ErrorInternalID {
		status = http.StatusTeapot
	}

	return status, err
}
