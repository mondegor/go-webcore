package observe

import (
	"bytes"
	"io"
)

type (
	// requestBody - декоратор io.ReadCloser.
	requestBody struct {
		originalBody io.ReadCloser // underlying reader
		size         int
		buf          bytes.Buffer
		bufSize      int
	}
)

// Read - реализует интерфейс io.ReadCloser.
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

// Close - реализует интерфейс io.ReadCloser.
func (rb *requestBody) Close() error {
	return rb.originalBody.Close()
}
