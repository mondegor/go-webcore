package observe

import (
	"bytes"
	"net/http"
)

type (
	// ResponseWriter - декоратор http.ResponseWriter для сбора статистики (статуса, кол-во записанных байт)
	// и возможности логирования записанных данных.
	ResponseWriter struct {
		http.ResponseWriter
		body       bytes.Buffer
		bufferSize int
		size       int
		statusCode int
	}
)

// NewResponseWriter - создаёт объект ResponseWriter.
func NewResponseWriter(w http.ResponseWriter, bufferSize int) *ResponseWriter {
	return &ResponseWriter{
		ResponseWriter: w,
		bufferSize:     bufferSize,
		statusCode:     http.StatusOK,
	}
}

// WriteHeader - устанавливает код ответа.
func (w *ResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

// Write - записывает переданные данные, подсчитывая их размер в байтах.
func (w *ResponseWriter) Write(buf []byte) (int, error) {
	n, err := w.ResponseWriter.Write(buf)
	if err != nil {
		return 0, err
	}

	w.size += n

	if w.bufferSize > 0 {
		if n > w.bufferSize {
			buf = buf[0:w.bufferSize]
		}

		w.bufferSize -= n // может стать отрицательным
		w.body.Write(buf)
	}

	return n, nil
}

// Content - возвращает копию переданных данных.
func (w *ResponseWriter) Content() []byte {
	return w.body.Bytes()
}

// Size - возвращает размер переданных данных (bytes).
func (w *ResponseWriter) Size() int {
	return w.size
}

// StatusCode - возвращает текущий код ответа.
func (w *ResponseWriter) StatusCode() int {
	return w.statusCode
}
