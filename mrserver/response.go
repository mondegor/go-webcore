package mrserver

import (
	"io"
	"net/http"

	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	ResponseEncoder interface {
		ContentType() string
		ContentTypeProblem() string
		Marshal(structure any) ([]byte, error)
	}

	ResponseSender interface {
		Send(w http.ResponseWriter, status int, structure any) error
		SendNoContent(w http.ResponseWriter) error
	}

	FileResponseSender interface {
		ResponseSender
		SendFile(w http.ResponseWriter, info mrtype.FileInfo, attachmentName string, file io.Reader) error
	}

	ErrorResponseSender interface {
		SendError(w http.ResponseWriter, r *http.Request, err error)
	}

	HttpErrorOverrideFunc func(err *mrerr.AppError) (int, *mrerr.AppError)
)
