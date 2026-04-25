package stat

import (
	"net/http"
	"strconv"
	"time"

	"github.com/mondegor/go-webcore/mrserver"
)

type (
	// RequestMetrics - сборщик метрик HTTP-запросов для системы мониторинга.
	// Собирает и отправляет:
	//  - Время выполнения запроса (гистограмма);
	//  - Размер полученного запроса (счётчик);
	//  - Размер отправленного ответа (счётчик).
	RequestMetrics struct {
		metrics mrserver.RequestObserve
	}
)

// NewRequestMetrics - создаёт сборщик метрик HTTP-запросов.
func NewRequestMetrics(metrics mrserver.RequestObserve) *RequestMetrics {
	return &RequestMetrics{
		metrics: metrics,
	}
}

// Enabled - всегда возвращает true, так как сбор метрик всегда включен.
func (rs *RequestMetrics) Enabled() bool {
	return true
}

// Emit - собирает и отправляет метрики HTTP-запроса.
// TODO: использовать нормализованный путь для группировки метрик по шаблону маршрута.
func (rs *RequestMetrics) Emit(r *http.Request, _ []byte, size int, _ []byte, responseSize int, duration time.Duration, status int) {
	method := r.Method
	path := r.URL.Path // TODO: из пути обрезать ID и другие параметры

	rs.metrics.SetStatusWithTime(method, path, strconv.Itoa(status), duration)
	rs.metrics.IncrementRequestSize(method, path, size)
	rs.metrics.IncrementResponseSize(method, path, responseSize)
}
