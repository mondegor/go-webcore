package mrserver

import (
	"net/http"
)

const (
	traceResponseBodyMaxLen = 1024
)

type (
	StatResponseWriter struct {
		http.ResponseWriter
		statusCode int
		bytes      int
		onWrite    func(buf []byte)
	}
)

func NewStatResponseWriter(w http.ResponseWriter, onWrite func(buf []byte)) *StatResponseWriter {
	return &StatResponseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
		onWrite:        onWrite,
	}
}

func (w *StatResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *StatResponseWriter) Write(buf []byte) (int, error) {
	bytes, err := w.ResponseWriter.Write(buf)
	w.bytes += bytes

	if len(buf) > traceResponseBodyMaxLen {
		buf = buf[0:traceResponseBodyMaxLen]
	}

	w.onWrite(buf)

	return bytes, err
}

func (w *StatResponseWriter) Bytes() int {
	return w.bytes
}

func (w *StatResponseWriter) StatusCode() int {
	return w.statusCode
}
