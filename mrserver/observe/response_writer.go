package observe

import (
	"bytes"
	"net/http"
)

type (
	// ResponseWriter - декоратор http.ResponseWriter для сбора статистики ответа.
	// Используется в ObserverHandler для мониторинга HTTP-трафика.
	//
	// Предоставляет:
	//  - Перехват HTTP-статуса ответа (по умолчанию 200 OK);
	//  - Подсчёт общего размера записанных данных;
	//  - Сохранение копии тела ответа (до указанного лимита) для логирования.
	ResponseWriter struct {
		http.ResponseWriter
		body       bytes.Buffer
		bufferSize int
		size       int
		statusCode int
	}
)

// NewResponseWriter - создаёт декорированный ResponseWriter для сбора статистики.
// Параметры:
//   - w - оригинальный http.ResponseWriter;
//   - bufferSize - максимальный размер копии тела ответа для логирования (0 = не сохранять).
func NewResponseWriter(w http.ResponseWriter, bufferSize int) *ResponseWriter {
	return &ResponseWriter{
		ResponseWriter: w,
		bufferSize:     bufferSize,
		statusCode:     http.StatusOK,
	}
}

// WriteHeader - перехватывает установку HTTP-кода ответа.
// Сохраняет код статуса для последующего чтения через StatusCode().
func (w *ResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

// Write - записывает данные в ответ, подсчитывая размер и сохраняя копию.
// Первые bufferSize байт сохраняются во внутренний буфер для логирования.
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

// Content - возвращает копию записанного тела ответа.
// Содержит не более bufferSize первых байт ответа.
func (w *ResponseWriter) Content() []byte {
	return w.body.Bytes()
}

// Size - возвращает общий размер записанных данных в байтах.
func (w *ResponseWriter) Size() int {
	return w.size
}

// StatusCode - возвращает установленный HTTP-код ответа.
func (w *ResponseWriter) StatusCode() int {
	return w.statusCode
}
