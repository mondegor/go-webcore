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

// NewStaticPeriodStrategy - создает стратегию статичного периода.
func NewStaticPeriodStrategy(value time.Duration) PeriodStrategy {
	return strategy.NewStaticPeriod(value)
}

// NewDispersionPeriodStrategy - создает стратегию периода со случайной дисперсией.
// Параметры:
//   - value - базовый период;
//   - ratio - коэффициент дисперсии [0..1), при 0 период статичен.
func NewDispersionPeriodStrategy(period time.Duration, ratio float64) PeriodStrategy {
	return strategy.NewDispersionPeriod(period, ratio)
}

// NewDelayedPeriodStrategy - создает стратегию периода с начальным запаздыванием.
// Параметры:
//   - delayed - начальная задержка перед выходом на номинальный период (может быть отрицательной);
//   - ratio - коэффициент дисперсии [0..1), добавляющий случайное отклонение к периоду;
//   - decay - коэффициент затухания задержки [0..1);
//   - periodStrategy - стратегия определения номинального периода.
func NewDelayedPeriodStrategy(
	delayed time.Duration,
	ratio float64,
	decay float64,
	periodStrategy PeriodStrategy,
) PeriodStrategy {
	return strategy.NewDelayedPeriod(delayed, ratio, decay, periodStrategy)
}

// NewDoubleDelayedStartStrategy - создает составную стратегию периода с удвоенным запаздыванием на старте.
// Комбинирует DelayedPeriod (с начальной задержкой, равной периоду, и с затуханием 50%)
// и DispersionPeriod (со случайной дисперсией ratio).
func NewDoubleDelayedStartStrategy(period time.Duration, ratio float64) PeriodStrategy {
	return strategy.NewDelayedPeriod(
		period,
		ratio,
		0.5, // затухание 50%
		strategy.NewDispersionPeriod(
			period,
			ratio,
		),
	)
}

// NewQuadQuickStartStrategy - создает составную стратегию периода с четырёх-кратным ускорением на старте.
// Комбинирует DelayedPeriod (с начальным ускорением в 4 раза больше периода, и с мгновенным затуханием)
// и DispersionPeriod (со случайной дисперсией ratio).
func NewQuadQuickStartStrategy(period time.Duration, ratio float64) PeriodStrategy {
	return strategy.NewDelayedPeriod(
		-(period/4)*3, // period / 4 = k * period + period -> k = - 3/4
		ratio,
		0, // мгновенное затухание
		strategy.NewDispersionPeriod(
			period,
			ratio,
		),
	)
}
