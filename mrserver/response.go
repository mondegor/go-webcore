package mrserver

import (
	"context"
	"io"
	"net/http"

	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrtype"
)

const (
	traceResponseMaxLen = 1024
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

	StatResponseWriter struct {
		http.ResponseWriter
		ctx        context.Context
		statusCode int
		bytes      int
	}
)

func NewStatResponseWriter(ctx context.Context, w http.ResponseWriter) *StatResponseWriter {
	return &StatResponseWriter{w, ctx, http.StatusOK, 0}
}

func (w *StatResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *StatResponseWriter) Write(buf []byte) (int, error) {
	bytes, err := w.ResponseWriter.Write(buf)
	w.bytes += bytes

	if len(buf) > traceResponseMaxLen {
		buf = buf[0:traceResponseMaxLen]
	}

	mrlog.Ctx(w.ctx).Trace().Bytes("response", buf).Msg("write response")

	return bytes, err
}
