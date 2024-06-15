package mrserver

import (
	"net/http"
)

const (
	traceResponseBodyMaxLen = 1024
)

type (
	// StatResponseWriter - comment struct.
	StatResponseWriter struct {
		http.ResponseWriter
		statusCode int
		bytes      int
		onWrite    func(buf []byte)
	}
)

// NewStatResponseWriter - создаёт объект StatResponseWriter.
func NewStatResponseWriter(w http.ResponseWriter, onWrite func(buf []byte)) *StatResponseWriter {
	return &StatResponseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
		onWrite:        onWrite,
	}
}

// WriteHeader - comment method.
func (w *StatResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

// Write - comment method.
func (w *StatResponseWriter) Write(buf []byte) (int, error) {
	bytes, err := w.ResponseWriter.Write(buf)
	w.bytes += bytes

	if len(buf) > traceResponseBodyMaxLen {
		buf = buf[0:traceResponseBodyMaxLen]
	}

	w.onWrite(buf)

	return bytes, err
}

// Bytes - comment method.
func (w *StatResponseWriter) Bytes() int {
	return w.bytes
}

// StatusCode - comment method.
func (w *StatResponseWriter) StatusCode() int {
	return w.statusCode
}
