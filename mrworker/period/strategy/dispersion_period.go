package strategy

import (
	"time"
)

type (
	// dispersionPeriod - стратегия периода со случайной дисперсией.
	// К базовому периоду добавляется случайное отклонение
	// в диапазоне [-ratio*period, +ratio*period].
	// Это решает проблему синхронизации параллельно запущенных тикеров,
	// предотвращая одновременное срабатывание множества задач,
	// что снижает пиковую нагрузку на систему.
	dispersionPeriod struct {
		period time.Duration
		ratio  float64
	}
)

// NewDispersionPeriod - создает стратегию периода со случайной дисперсией.
// Параметры:
//   - value - базовый период;
//   - ratio - коэффициент дисперсии [0..1), при 0 период статичен.
func NewDispersionPeriod(period time.Duration, ratio float64) Strategy {
	return &dispersionPeriod{
		period: fixedPeriod(period),
		ratio:  fixedCoefficient(ratio),
	}
}

// Period - возвращает период со случайным отклонением.
// При ratio = 0 возвращается базовый период.
// При ratio > 0 к периоду добавляется случайное значение в диапазоне [-ratio*period, +ratio*period].
func (p *dispersionPeriod) Period() time.Duration {
	return calcPeriod(p.period, p.ratio)
}
