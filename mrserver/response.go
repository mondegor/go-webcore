package mrserver

import (
	"context"
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
		SendBytes(w http.ResponseWriter, status int, body []byte) error
		SendNoContent(w http.ResponseWriter) error
	}

	FileResponseSender interface {
		ResponseSender
		SendFile(ctx context.Context, w http.ResponseWriter, file mrtype.File) error
		SendAttachmentFile(ctx context.Context, w http.ResponseWriter, file mrtype.File) error
	}

	ErrorResponseSender interface {
		SendError(w http.ResponseWriter, r *http.Request, err error)
	}

	HTTPErrorOverrideFunc func(err *mrerr.AppError) (int, *mrerr.AppError)
)
