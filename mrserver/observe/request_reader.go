package observe

import (
	"net/http"
)

type (
	// RequestReader - псевдо декоратор http.Request для сбора статистики (кол-во прочитанных байт)
	// и возможности логирования считанных данных.
	RequestReader struct {
		request  *http.Request
		headSize int
		body     *requestBody
	}
)

// NewRequestReader - создаёт объект RequestReader.
// WARNING: the Body of the original http.Request can be replaced.
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

// Request - возвращает оригинальный запрос.
func (r *RequestReader) Request() *http.Request {
	return r.request
}

// Content - возвращает копию прочитанных данных.
func (r *RequestReader) Content() []byte {
	if r.body == nil {
		return nil
	}

	return r.body.buf.Bytes()
}

// Size - возвращает размер считанных данных (bytes).
func (r *RequestReader) Size() int {
	if r.body == nil {
		return r.headSize
	}

	return r.headSize + r.body.size
}
