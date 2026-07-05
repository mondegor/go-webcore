package observe

import (
	"net/http"
)

type (
	// RequestReader - обёртка над http.Request для сбора статистики запроса.
	//
	// Предоставляет:
	//  - Подсчёт общего размера прочитанных данных (заголовок + тело);
	//  - Сохранение копии тела запроса (до указанного лимита) для логирования.
	RequestReader struct {
		request  *http.Request
		headSize int
		body     *requestBody
	}
)

// NewRequestReader - создаёт обёртку для сбора статистики HTTP-запроса.
//
// Важно:
//   - При bufferSize > 0 и методе POST/PUT заменяет r.Body на декорированную версию;
//   - Для методов без тела (GET, DELETE и т.д.) тело не перехватывается;
func NewRequestReader(r *http.Request, bufferSize int) *RequestReader {
	var body *requestBody

	if bufferSize > 0 {
		if r.Method == http.MethodPost || r.Method == http.MethodPut {
			body = &requestBody{
				originalBody: r.Body,
				bufSize:      bufferSize,
			}

			r.Body = body
		}
	}

	return &RequestReader{
		request:  r,
		headSize: len(r.URL.RawQuery),
		body:     body,
	}
}

// Request - возвращает оригинальный HTTP-запрос.
// Тело запроса может быть заменено на декорированную версию для отслеживания.
func (r *RequestReader) Request() *http.Request {
	return r.request
}

// Content - возвращает копию прочитанного тела запроса.
// Содержит не более bufferSize первых байт тела.
func (r *RequestReader) Content() []byte {
	if r.body == nil {
		return nil
	}

	return r.body.buf.Bytes()
}

// Size - возвращает общий размер запроса в байтах (длина RawQuery + размер тела).
func (r *RequestReader) Size() int {
	if r.body == nil {
		return r.headSize
	}

	return r.headSize + r.body.size
}
