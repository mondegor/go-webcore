package mrserver

import (
	"io"
	"net/http"
)

type (
	// StatRequest - comment struct.
	StatRequest struct {
		request *http.Request
		body    *requestBody
	}

	requestBody struct {
		b     io.ReadCloser // underlying reader
		bytes int
	}
)

// NewStatRequest - создаёт объект StatRequest.
// WARNING: the Body of the original http.Request will be replaced.
func NewStatRequest(r *http.Request) *StatRequest {
	body := &requestBody{
		b:     r.Body,
		bytes: len(r.URL.RawQuery),
	}

	r.Body = body

	return &StatRequest{
		request: r,
		body:    body,
	}
}

// Request - comment method.
func (r *StatRequest) Request() *http.Request {
	return r.request
}

// Bytes - comment method.
func (r *StatRequest) Bytes() int {
	return r.body.bytes
}

// Read - comment method.
func (rb *requestBody) Read(p []byte) (n int, err error) {
	rb.bytes += len(p)

	return rb.b.Read(p)
}

// Close - comment method.
func (rb *requestBody) Close() error {
	return rb.b.Close()
}
