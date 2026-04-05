package stat

import (
	"net/http"
	"strconv"
	"time"
)

type (
	// RequestMetrics - собирает и отправляет метрики HTTP-запросов.
	RequestMetrics struct {
		metrics requestMetrics
	}

	// requestMetrics - внутренний интерфейс для отправки метрик в систему мониторинга.
	requestMetrics interface {
		SetStatusWithTime(method, location, status string, duration time.Duration)
		IncrementRequestSize(method, location string, size int)
		IncrementResponseSize(method, location string, size int)
	}
)

// NewRequestMetrics - создаёт объект RequestMetrics.
func NewRequestMetrics(metrics requestMetrics) *RequestMetrics {
	return &RequestMetrics{
		metrics: metrics,
	}
}

// Enabled - возвращает true, так как сбор метрик всегда включен.
func (rs *RequestMetrics) Enabled() bool {
	return true
}

// Emit - собирает и отправляет метрики HTTP запроса.
func (rs *RequestMetrics) Emit(r *http.Request, _ []byte, size int, _ []byte, responseSize int, duration time.Duration, status int) {
	method := r.Method
	path := r.URL.Path // TODO: из пути обрезать ID и другие параметры

	rs.metrics.SetStatusWithTime(method, path, strconv.Itoa(status), duration)
	rs.metrics.IncrementRequestSize(method, path, size)
	rs.metrics.IncrementResponseSize(method, path, responseSize)
}
