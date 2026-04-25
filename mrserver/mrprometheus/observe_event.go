package mrprometheus

import (
	"github.com/prometheus/client_golang/prometheus"
)

type (
	// ObserveEvent - сборщик метрик событий от источников данных.
	// Используется для мониторинга и алертинга на основе Prometheus.
	ObserveEvent struct {
		eventCount *prometheus.CounterVec
	}
)

// NewObserveEvent - создаёт сборщик метрик событий.
// Параметры:
//   - namespace - пространство имён метрик (например: имя приложения);
//   - subsystem - подсистема (например: "events").
func NewObserveEvent(namespace, subsystem string) *ObserveEvent {
	return &ObserveEvent{
		eventCount: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "event_sent_count",
				Help:      "Cumulative count of events for data sources.",
			},
			[]string{"event", "source"},
		),
	}
}

// Collectors - возвращает срез всех коллекторов метрик событий для регистрации в Prometheus.
// Используется при инициализации реестра метрик.
func (o *ObserveEvent) Collectors() []prometheus.Collector {
	return []prometheus.Collector{
		o.eventCount,
	}
}

// IncrementEvent - инкрементирует счётчик указанного события от источника данных.
func (o *ObserveEvent) IncrementEvent(event, source string) {
	o.eventCount.With(prometheus.Labels{"event": event, "source": source}).Inc()
}
