package mrinit

import (
	"github.com/prometheus/client_golang/prometheus"
)

type (
	// Prometheus - управляет реестром метрик Prometheus и коллекторами.
	Prometheus struct {
		registry *prometheus.Registry
		list     []prometheus.Collector
	}
)

// NewPrometheus - создаёт объект Prometheus.
func NewPrometheus() *Prometheus {
	return &Prometheus{
		registry: prometheus.NewRegistry(),
	}
}

// Registry - возвращает реестр метрик Prometheus.
func (rl *Prometheus) Registry() *prometheus.Registry {
	return rl.registry
}

// Add - добавляет коллектор в список для последующей регистрации.
func (rl *Prometheus) Add(collector ...prometheus.Collector) {
	rl.list = append(rl.list, collector...)
}

// Register - регистрирует все добавленные коллекторы в реестре Prometheus.
func (rl *Prometheus) Register() error {
	for _, collector := range rl.list {
		if err := rl.registry.Register(collector); err != nil {
			return err
		}
	}

	return nil
}
