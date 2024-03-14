package mrserver

import (
	"context"
	"net/http"

	"github.com/mondegor/go-webcore/mrlog"
)

const (
	traceResponseMaxLen = 1024
)

type (
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
