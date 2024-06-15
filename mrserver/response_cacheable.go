package mrserver

import (
	"net/http"
)

type (
	// CacheableResponseWriter - comment struct.
	CacheableResponseWriter struct {
		http.ResponseWriter
		statusCode int
		body       []byte
	}
)

// NewCacheableResponseWriter - создаёт объект CacheableResponseWriter.
func NewCacheableResponseWriter(w http.ResponseWriter) *CacheableResponseWriter {
	return &CacheableResponseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
		body:           nil,
	}
}

// WriteHeader - comment method.
func (w *CacheableResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

// Write - comment method.
func (w *CacheableResponseWriter) Write(buf []byte) (int, error) {
	if len(w.body) > 0 {
		w.body = append(w.body, buf...)
	} else {
		w.body = buf
	}

	return w.ResponseWriter.Write(buf)
}

// StatusCode - comment method.
func (w *CacheableResponseWriter) StatusCode() int {
	return w.statusCode
}

// Body - comment method.
func (w *CacheableResponseWriter) Body() []byte {
	return w.body
}
