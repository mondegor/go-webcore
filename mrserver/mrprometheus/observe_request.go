package mrprometheus

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type (
	// ObserveRequest - метрики сетевых запросов.
	ObserveRequest struct {
		requestDuration *prometheus.HistogramVec
		requestSize     *prometheus.CounterVec
		responseSize    *prometheus.CounterVec
	}
)

// NewObserveRequest - создаёт объект ObserveRequest.
func NewObserveRequest(namespace, subsystem string) *ObserveRequest {
	return &ObserveRequest{
		requestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "request_duration_seconds",
				Help:      "Request executed time.",
				Buckets:   []float64{0.001, 0.005, 0.025, 0.2, 1, 4, 8},
			},
			[]string{"method", "location", "status"},
		),
		requestSize: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "request_received_bytes_total",
				Help:      "Size in bytes of received information.",
			},
			[]string{"method", "location"},
		),
		responseSize: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "response_sent_bytes_total",
				Help:      "Size in bytes of sent information.",
			},
			[]string{"method", "location"},
		),
	}
}

// Collectors - возвращает всех собирателей метрик сетевых запросов.
func (o *ObserveRequest) Collectors() []prometheus.Collector {
	return []prometheus.Collector{
		o.requestDuration,
		o.requestSize,
		o.responseSize,
	}
}

// SetStatusWithTime - устанавливает статус запроса (200, 404, 500, и т.д) и время исполнения запроса для указанного URL.
func (o *ObserveRequest) SetStatusWithTime(method, location, status string, duration time.Duration) {
	o.requestDuration.With(prometheus.Labels{"method": method, "location": location, "status": status}).Observe(float64(duration.Microseconds()))
}

// IncrementRequestSize - добавляет размер тела запроса в байтах для указанного URL.
func (o *ObserveRequest) IncrementRequestSize(method, location string, size int) {
	o.requestSize.With(prometheus.Labels{"method": method, "location": location}).Add(float64(size))
}

// IncrementResponseSize - добавляет размер тела ответа в байтах для указанного URL.
func (o *ObserveRequest) IncrementResponseSize(method, location string, size int) {
	o.responseSize.With(prometheus.Labels{"method": method, "location": location}).Add(float64(size))
}
