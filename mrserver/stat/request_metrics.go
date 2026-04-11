package stat

import (
	"net/http"
	"strconv"
	"time"
)

type (
	// RequestMetrics - сборщик метрик HTTP-запросов для системы мониторинга.
	//
	// Собирает и отправляет:
	//  - Время выполнения запроса (гистограмма)
	//  - Размер полученного запроса (счётчик)
	//  - Размер отправленного ответа (счётчик)
	//
	// Реализует интерфейс mrserver.RequestStat для использования в ObserverHandler.
	RequestMetrics struct {
		// metrics - внутренний отправитель метрик (например: mrprometheus.ObserveRequest).
		metrics requestMetrics
	}

	// requestMetrics - внутренний интерфейс для отправки метрик в систему мониторинга.
	// Реализуется, например, mrprometheus.ObserveRequest.
	requestMetrics interface {
		SetStatusWithTime(method, location, status string, duration time.Duration)
		IncrementRequestSize(method, location string, size int)
		IncrementResponseSize(method, location string, size int)
	}
)

// NewRequestMetrics - создаёт сборщик метрик HTTP-запросов.
func NewRequestMetrics(metrics requestMetrics) *RequestMetrics {
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
