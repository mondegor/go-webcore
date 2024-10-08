package mrprometheus

import (
	"github.com/prometheus/client_golang/prometheus"
)

type (
	// ObserveEvent - метрики событий источников данных.
	ObserveEvent struct {
		eventCount *prometheus.CounterVec
	}
)

// NewObserveEvent - создаёт объект ObserveEvent.
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

// Collectors - возвращает всех собирателей метрик событий источников данных.
func (o *ObserveEvent) Collectors() []prometheus.Collector {
	return []prometheus.Collector{
		o.eventCount,
	}
}

// IncrementEvent - увеличивает счётчик указанного события для указанного источника данных.
func (o *ObserveEvent) IncrementEvent(event, source string) {
	o.eventCount.With(prometheus.Labels{"event": event, "source": source}).Add(1)
}
