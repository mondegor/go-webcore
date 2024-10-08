package mrserver

import (
	"bytes"
	"net/http"
)

type (
	// CacheableResponseWriter - comment struct.
	CacheableResponseWriter struct {
		http.ResponseWriter
		statusCode int
		body       bytes.Buffer
	}
)

// NewCacheableResponseWriter - создаёт объект CacheableResponseWriter.
func NewCacheableResponseWriter(w http.ResponseWriter) *CacheableResponseWriter {
	return &CacheableResponseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}
}

// WriteHeader - comment method.
func (w *CacheableResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

// Write - записывает переданные данные.
func (w *CacheableResponseWriter) Write(buf []byte) (int, error) {
	n, err := w.ResponseWriter.Write(buf)
	if err != nil {
		return n, err
	}

	if n > 0 {
		w.body.Write(buf[0:n])
	}

	return n, nil
}

// StatusCode - возвращает текущий код ответа.
func (w *CacheableResponseWriter) StatusCode() int {
	return w.statusCode
}

// Content - возвращает копию переданных данных.
func (w *CacheableResponseWriter) Content() []byte {
	return w.body.Bytes()
}
