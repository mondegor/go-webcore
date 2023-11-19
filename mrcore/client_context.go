package mrcore

import (
	"context"
	"io"
	"net/http"

	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	ClientContext interface {
		Request() *http.Request
		ParamFromPath(name string) string

		Context() context.Context
		WithContext(ctx context.Context) ClientContext

		Writer() http.ResponseWriter

		Validate(structRequest any) error

		SendResponse(status int, structResponse any) error
		SendResponseNoContent() error
		SendFile(info mrtype.FileInfo, attachmentName string, file io.Reader) error
	}

	ClientErrorWrapperFunc func(err *mrerr.AppError) (int, *mrerr.AppError)
)
