package observe

import (
	"bytes"
	"io"
)

type (
	// requestBody - декоратор io.ReadCloser для подсчёта прочитанных байт
	// и сохранения копии тела запроса (до указанного лимита).
	// Используется внутри RequestReader для сбора статистики и логирования.
	// При чтении пропускает данные через себя, подсчитывая общий размер
	// и записывая копию первых bufSize байт во внутренний буфер.
	requestBody struct {
		originalBody io.ReadCloser
		size         int
		buf          bytes.Buffer
		bufSize      int
	}
)

// Read - читает данные из исходного тела запроса.
// Подсчитывает общий размер прочитанных байт и сохраняет копию
// первых bufSize байт во внутренний буфер для последующего логирования.
func (rb *requestBody) Read(p []byte) (n int, err error) {
	n, err = rb.originalBody.Read(p)
	rb.size += n

	if rb.bufSize > 0 && n > 0 {
		if n > rb.bufSize {
			n = rb.bufSize
		}

		rb.bufSize -= n
		rb.buf.Write(p[0:n])
	}

	return n, err
}

// Close - закрывает исходный читатель тела запроса.
func (rb *requestBody) Close() error {
	return rb.originalBody.Close()
}
