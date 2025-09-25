package mrinit

import (
	"github.com/prometheus/client_golang/prometheus"
)

type (
	// Prometheus - comment struct.
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

// Registry - comment method.
func (rl *Prometheus) Registry() *prometheus.Registry {
	return rl.registry
}

// Add - comment method.
func (rl *Prometheus) Add(collector ...prometheus.Collector) {
	rl.list = append(rl.list, collector...)
}

// Register - comment method.
func (rl *Prometheus) Register() error {
	for _, collector := range rl.list {
		if err := rl.registry.Register(collector); err != nil {
			return err
		}
	}

	return nil
}
