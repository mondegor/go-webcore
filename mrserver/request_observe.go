package mrserver

import (
	"time"
)

type (
	// RequestObserve - интерфейс для отправки метрик в систему мониторинга.
	RequestObserve interface {
		SetStatusWithTime(method, location, status string, duration time.Duration)
		IncrementRequestSize(method, location string, size int)
		IncrementResponseSize(method, location string, size int)
	}

	// nopRequestObserve - заглушка, реализующая интерфейс RequestObserve.
	// Игнорирует передаваемые метрики.
	nopRequestObserve struct{}
)

// NopRequestObserve - создаёт RequestObserve, который игнорирует все метрики.
func NopRequestObserve() RequestObserve {
	return nopRequestObserve{}
}

// SetStatusWithTime - имитирует сбор метрик.
func (o nopRequestObserve) SetStatusWithTime(_, _, _ string, _ time.Duration) {}

// IncrementRequestSize - имитирует сбор метрик.
func (o nopRequestObserve) IncrementRequestSize(_, _ string, _ int) {}

// IncrementResponseSize - имитирует сбор метрик.
func (o nopRequestObserve) IncrementResponseSize(_, _ string, _ int) {}
