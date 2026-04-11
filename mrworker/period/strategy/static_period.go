package strategy

import (
	"time"
)

type (
	// staticPeriod - стратегия фиксированного (статичного) периода.
	// Всегда возвращает одно и то же значение периода, заданное при создании.
	staticPeriod time.Duration
)

// NewStaticPeriod - создает стратегию статичного периода.
func NewStaticPeriod(period time.Duration) Strategy {
	return staticPeriod(fixedPeriod(period))
}

// Period - возвращает текущий период.
func (p staticPeriod) Period() time.Duration {
	return time.Duration(p)
}
