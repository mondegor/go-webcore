package mrprometheus

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type (
	// ObserveRequest - сборщик метрик HTTP-запросов.
	// Используется для мониторинга производительности и трафика HTTP-сервера.
	//
	// Предоставляет:
	//  - requestDuration - гистограмма времени выполнения запросов;
	//  - requestSize - счётчик размера полученных данных (байты);
	//  - responseSize - счётчик размера отправленных данных (байты).
	ObserveRequest struct {
		requestDuration *prometheus.HistogramVec
		requestSize     *prometheus.CounterVec
		responseSize    *prometheus.CounterVec
	}
)

// NewObserveRequest - создаёт сборщик метрик HTTP-запросов.
// Параметры:
//   - namespace - пространство имён метрик (например: имя приложения);
//   - subsystem - подсистема (например: "http").
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

// Collectors - возвращает срез всех коллекторов метрик запросов для регистрации в Prometheus.
// Используется при инициализации реестра метрик.
func (o *ObserveRequest) Collectors() []prometheus.Collector {
	return []prometheus.Collector{
		o.requestDuration,
		o.requestSize,
		o.responseSize,
	}
}

// SetStatusWithTime - записывает длительность выполнения HTTP-запроса с указанием статуса.
// Параметры:
//   - method - HTTP-метод (GET, POST, и т.д.);
//   - location - путь запроса (URL-шаблон);
//   - status - HTTP-код ответа (например: "200", "404", "500");
//   - duration - время выполнения запроса.
func (o *ObserveRequest) SetStatusWithTime(method, location, status string, duration time.Duration) {
	o.requestDuration.With(prometheus.Labels{"method": method, "location": location, "status": status}).Observe(float64(duration.Microseconds()))
}

// IncrementRequestSize - добавляет размер тела полученного HTTP-запроса в байтах.
// Параметры:
//   - method - HTTP-метод запроса;
//   - location - путь запроса (URL-шаблон);
//   - size - размер тела запроса в байтах.
func (o *ObserveRequest) IncrementRequestSize(method, location string, size int) {
	o.requestSize.With(prometheus.Labels{"method": method, "location": location}).Add(float64(size))
}

// IncrementResponseSize - добавляет размер тела отправленного HTTP-ответа в байтах.
// Параметры:
//   - method - HTTP-метод запроса;
//   - location - путь запроса (URL-шаблон);
//   - size - размер тела ответа в байтах.
func (o *ObserveRequest) IncrementResponseSize(method, location string, size int) {
	o.responseSize.With(prometheus.Labels{"method": method, "location": location}).Add(float64(size))
}
