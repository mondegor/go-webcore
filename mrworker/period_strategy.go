package mrworker

import (
	"time"

	"github.com/mondegor/go-webcore/mrworker/period/strategy"
)

type (
	// PeriodStrategy - интерфейс стратегии определения периода тиков.
	PeriodStrategy interface {
		Period() time.Duration
	}
)

// NewStaticPeriod - создает стратегию статичного периода.
func NewStaticPeriod(value time.Duration) PeriodStrategy {
	return strategy.NewStaticPeriod(value)
}

// NewDispersionPeriod - создает стратегию периода со случайной дисперсией.
// Параметры:
//   - value - базовый период;
//   - ratio - коэффициент дисперсии [0..1), при 0 период статичен.
func NewDispersionPeriod(period time.Duration, ratio float64) PeriodStrategy {
	return strategy.NewDispersionPeriod(period, ratio)
}

// NewDelayedPeriod - создает стратегию периода с начальным запаздыванием.
// Параметры:
//   - delayed - начальная задержка перед выходом на номинальный период (может быть отрицательной);
//   - ratio - коэффициент дисперсии [0..1), добавляющий случайное отклонение к периоду;
//   - decay - коэффициент затухания задержки [0..1);
//   - periodStrategy - стратегия определения номинального периода.
func NewDelayedPeriod(
	delayed time.Duration,
	ratio float64,
	decay float64,
	periodStrategy PeriodStrategy,
) PeriodStrategy {
	return strategy.NewDelayedPeriod(delayed, ratio, decay, periodStrategy)
}

// NewDoubleDelayedPeriod - создает составную стратегию периода с удвоенным запаздыванием на старте.
// Комбинирует DelayedPeriod (с начальной задержкой, равной периоду, и затуханием 0.5)
// и DispersionPeriod (со случайной дисперсией для предотвращения синхронизации тикеров).
func NewDoubleDelayedPeriod(period time.Duration, ratio float64) PeriodStrategy {
	return strategy.NewDelayedPeriod(
		period,
		ratio,
		0.5,
		strategy.NewDispersionPeriod(
			period,
			ratio,
		),
	)
}
