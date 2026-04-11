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

// NewPrometheus - создаёт и инициализирует объект Prometheus с новым реестром метрик.
// Реестр изолирован от глобального реестра Prometheus для лучшей контролируемости.
func NewPrometheus() *Prometheus {
	return &Prometheus{
		registry: prometheus.NewRegistry(),
	}
}

// Registry - возвращает внутренний реестр метрик Prometheus.
// Реестр используется для регистрации и хранения всех метрик приложения.
// Может быть передан HTTP-серверу для exposition метрик.
func (rl *Prometheus) Registry() *prometheus.Registry {
	return rl.registry
}

// Add - добавляет коллекторы метрик во внутренний список для последующей регистрации.
// Коллекторы не регистрируются сразу, а накапливаются для пакетной регистрации через Register().
func (rl *Prometheus) Add(collector ...prometheus.Collector) {
	rl.list = append(rl.list, collector...)
}

// Register - регистрирует все добавленные через Add() коллекторы в реестре Prometheus.
// Вызывается после того, как все коллекторы были добавлены.
func (rl *Prometheus) Register() error {
	for _, collector := range rl.list {
		if err := rl.registry.Register(collector); err != nil {
			return err
		}
	}

	return nil
}
