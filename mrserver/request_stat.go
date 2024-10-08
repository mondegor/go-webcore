package mrserver

import (
	"net/http"
)

type (
	// StatRequestReader - псевдо декоратор http.Request для сбора статистики (кол-во прочитанных байт)
	// и возможности логирования считанных данных.
	StatRequestReader struct {
		request  *http.Request
		headSize int
		body     *requestBody
	}
)

// NewStatRequestReader - создаёт объект StatRequestReader.
// WARNING: the Body of the original http.Request can be replaced.
func NewStatRequestReader(r *http.Request, bufferSize int) *StatRequestReader {
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

	return &StatRequestReader{
		request:  r,
		headSize: len(r.URL.RawQuery),
		body:     body,
	}
}

// Request - возвращает оригинальный запрос.
func (r *StatRequestReader) Request() *http.Request {
	return r.request
}

// Size - возвращает размер считанных данных (bytes).
func (r *StatRequestReader) Size() int {
	if r.body == nil {
		return r.headSize
	}

	return r.headSize + r.body.size
}

// HasContent - имеется ли у запроса тело.
func (r *StatRequestReader) HasContent() bool {
	return r.body != nil
}

// Content - возвращает копию прочитанных данных.
func (r *StatRequestReader) Content() []byte {
	if r.body == nil {
		return nil
	}

	return r.body.buf.Bytes()
}
