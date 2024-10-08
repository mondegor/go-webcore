package mrserver

import (
	"bytes"
	"net/http"
)

type (
	// StatResponseWriter - декоратор http.ResponseWriter для сбора статистики (статуса, кол-во записанных байт)
	// и возможности логирования записанных данных.
	StatResponseWriter struct {
		http.ResponseWriter
		statusCode int
		size       int
		body       bytes.Buffer
		bufSize    int
	}
)

// NewStatResponseWriter - создаёт объект StatResponseWriter.
func NewStatResponseWriter(w http.ResponseWriter, bufferSize int) *StatResponseWriter {
	return &StatResponseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
		bufSize:        bufferSize,
	}
}

// WriteHeader - устанавливает код ответа.
func (w *StatResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

// Write - записывает переданные данные, подсчитывая их размер в байтах.
func (w *StatResponseWriter) Write(buf []byte) (int, error) {
	n, err := w.ResponseWriter.Write(buf)
	w.size += n

	if w.bufSize > 0 {
		if n > w.bufSize {
			buf = buf[0:w.bufSize]
		}

		w.bufSize -= n
		w.body.Write(buf[0:n])
	}

	return n, err
}

// StatusCode - возвращает текущий код ответа.
func (w *StatResponseWriter) StatusCode() int {
	return w.statusCode
}

// Size - возвращает размер переданных данных (bytes).
func (w *StatResponseWriter) Size() int {
	return w.size
}

// Content - возвращает копию переданных данных.
func (w *StatResponseWriter) Content() []byte {
	return w.body.Bytes()
}
