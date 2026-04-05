package mrserver

import (
	"net/http"
	"time"
)

type (
	// RequestStat - определяет интерфейс для сбора
	// и отправки статистики HTTP-запросов.
	RequestStat interface {
		Enabled() bool
		Emit(r *http.Request, body []byte, size int, responseBody []byte, responseSize int, duration time.Duration, status int)
	}
)
