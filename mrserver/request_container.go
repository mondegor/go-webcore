package mrserver

import (
	"net/http"
	"time"
)

type (
	// RequestContainer - comment struct.
	RequestContainer struct {
		list []RequestStat
	}
)

// NewRequestContainer - создаёт объект RequestContainer.
func NewRequestContainer(list ...RequestStat) *RequestContainer {
	n := 0

	for _, item := range list {
		if item.Enabled() {
			n++
		}
	}

	if n == 0 {
		return &RequestContainer{}
	}

	if n > 0 {
		j := 0

		for i := 0; i < len(list); i++ {
			if list[i].Enabled() {
				list[j] = list[i]
				j++
			}
		}

		list = list[:j:j]
	}

	return &RequestContainer{
		list: list,
	}
}

// Enabled - comment method.
func (rs *RequestContainer) Enabled() bool {
	return len(rs.list) > 0
}

// Emit - функция трассировки http запроса.
func (rs *RequestContainer) Emit(r *http.Request, body []byte, size int, responseBody []byte, responseSize int, duration time.Duration, status int) {
	for i := range rs.list {
		rs.list[i].Emit(r, body, size, responseBody, responseSize, duration, status)
	}
}
