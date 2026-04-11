package strategy

import (
	"sync/atomic"
	"time"
)

type (
	// delayedPeriod - стратегия периода с начальным запаздыванием, которое
	// постепенно затухает. При каждом вызове Period() начальная задержка/ускорение
	// уменьшается коэффициентом decay, пока не станет равной базовому периоду.
	delayedPeriod struct {
		periodStrategy Strategy
		ratio          float64
		decay          float64
		delayed        atomic.Int64 // nanoseconds
	}
)

// NewDelayedPeriod - создает стратегию периода с начальным запаздыванием.
// Параметры:
//   - delayed - начальная задержка перед выходом на номинальный период (если отрицательна - ускорение);
//   - ratio - коэффициент дисперсии [0..1), добавляющий случайное отклонение к периоду;
//   - decay - коэффициент затухания задержки [0..1);
//   - periodStrategy - стратегия определения номинального периода.
func NewDelayedPeriod(
	delayed time.Duration,
	ratio float64,
	decay float64,
	periodStrategy Strategy,
) Strategy {
	if periodStrategy == nil {
		periodStrategy = NewStaticPeriod(minTimePeriod)
	}

	p := &delayedPeriod{
		periodStrategy: periodStrategy,
		ratio:          fixedCoefficient(ratio),
		decay:          fixedCoefficient(decay),
	}

	p.delayed.Store(p.truncateMs(delayed))

	return p
}

// Period - возвращает текущий период с учетом затухающей задержки/ускорения и дисперсии.
// При первом вызове к базовому периоду добавляется полная начальная задержка/ускорение.
// При каждом последующем вызове задержка/ускорение уменьшается коэффициентом decay.
// Когда задержка/ускорение затухает до нуля, возвращается только период определённый стратегией.
func (p *delayedPeriod) Period() time.Duration {
	delayed := p.delayed.Load()

	// если задержка/ускорение уже затухла, то сразу возвращается только период
	if delayed == 0 {
		return p.periodStrategy.Period()
	}

	// вычисляется задержка/ускорение для следующего шага с точностью до секунды
	newDelayed := calcPeriod(
		time.Duration(float64(delayed)*p.decay), // если отрицательна - ускорение
		p.ratio,
	)

	p.delayed.Store(p.truncateMs(newDelayed)) // best-effort

	return fixedPeriod(time.Duration(delayed) + p.periodStrategy.Period())
}

func (p *delayedPeriod) truncateMs(value time.Duration) int64 {
	return value.Milliseconds() * 1e6
}
