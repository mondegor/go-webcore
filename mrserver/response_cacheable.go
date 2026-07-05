package mrserver

import (
	"bytes"
	"net/http"
)

type (
	// CacheableResponseWriter - декоратор http.ResponseWriter для перехвата ответа.
	// Используется для кэширования HTTP-ответов (например: в middleware идемпотентности).
	// Записывает ответ одновременно в оригинальный ResponseWriter и во внутренний буфер,
	// позволяя впоследствии получить полную копию ответа для сохранения в кэш.
	CacheableResponseWriter struct {
		http.ResponseWriter
		statusCode int
		body       bytes.Buffer
	}
)

// NewCacheableResponseWriter - создаёт декоратор для перехвата HTTP-ответа.
func NewCacheableResponseWriter(w http.ResponseWriter) *CacheableResponseWriter {
	return &CacheableResponseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}
}

// WriteHeader - перехватывает и записывает HTTP-код ответа.
func (w *CacheableResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

// Write - записывает данные в ответ и сохраняет копию во внутренний буфер.
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

// StatusCode - возвращает установленный HTTP-код ответа.
func (w *CacheableResponseWriter) StatusCode() int {
	return w.statusCode
}

// Content - возвращает полную копию записанных данных ответа.
func (w *CacheableResponseWriter) Content() []byte {
	return w.body.Bytes()
}
