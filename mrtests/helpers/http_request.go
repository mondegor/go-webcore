package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
)

type (
	// HttpRequest - помощник для работы с http запросом.
	HttpRequest struct {
		r   *http.Request
		err error
	}
)

// NewHttpRequest - создаёт объект HttpRequest.
func NewHttpRequest(method, target string, structRequest any) *HttpRequest {
	requestBody, err := json.Marshal(structRequest)
	if err != nil {
		return &HttpRequest{
			err: fmt.Errorf("marshal failed: %w", err),
		}
	}

	return &HttpRequest{
		r: httptest.NewRequest(method, target, bytes.NewBuffer(requestBody)),
	}
}

// Request - возвращает http.Request.
func (r *HttpRequest) Request() *http.Request {
	return r.r
}

// Exec - запускает указанный обработчик и возвращает результат.
func (r *HttpRequest) Exec(handler http.Handler, structResponse any) (statusCode int, err error) {
	if r.err != nil {
		return 0, r.err
	}

	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r.r)

	response := w.Result()
	defer response.Body.Close()

	if err = json.Unmarshal(w.Body.Bytes(), &structResponse); err != nil {
		return 0, err
	}

	return response.StatusCode, nil
}
