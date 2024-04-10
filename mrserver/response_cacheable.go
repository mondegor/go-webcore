package mrserver

import (
	"net/http"
)

type (
	CacheableResponseWriter struct {
		http.ResponseWriter
		statusCode int
		body       []byte
	}
)

func NewCacheableResponseWriter(w http.ResponseWriter) *CacheableResponseWriter {
	return &CacheableResponseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
		body:           []byte{},
	}
}

func (w *CacheableResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *CacheableResponseWriter) Write(buf []byte) (int, error) {
	if len(w.body) > 0 {
		w.body = append(w.body, buf...)
	} else {
		w.body = buf
	}

	return w.ResponseWriter.Write(buf)
}

func (w *CacheableResponseWriter) StatusCode() int {
	return w.statusCode
}

func (w *CacheableResponseWriter) Body() []byte {
	return w.body
}
