package strategy

import (
	"math"
	"math/rand/v2"
	"time"
)

const (
	// minTimePeriod - минимально допустимый период тика.
	minTimePeriod = time.Millisecond
)

type (
	// Strategy - интерфейс стратегии определения периода тиков.
	// Метод Period() возвращает актуальную длительность периода,
	// которая может меняться в зависимости от реализации стратегии.
	Strategy interface {
		Period() time.Duration
	}
)

// calcPeriod - вычисляет период со случайным отклонением.
// Генерируется случайное значение в диапазоне
// [value-value*ratio, value+value*ratio] с равномерным распределением.
func calcPeriod(value time.Duration, ratio float64) time.Duration {
	delta := time.Duration(float64(value) * ratio)
	minInterval := value - delta
	maxInterval := value + delta

	if minInterval > maxInterval {
		minInterval, maxInterval = maxInterval, minInterval
	}

	// overflow guard
	if maxInterval >= math.MaxInt64-minInterval {
		minInterval, maxInterval = 0, time.Minute
	}

	return minInterval + time.Duration(rand.Int64N(int64(maxInterval-minInterval)+1)) //nolint:gosec
}

// fixedPeriod - предотвращает нулевые и отрицательные значения.
func fixedPeriod(value time.Duration) time.Duration {
	if value <= 0 {
		return minTimePeriod
	}

	return value
}

// fixedCoefficient - предотвращает значения больше или равные 1 и отрицательные значения.
func fixedCoefficient(value float64) float64 {
	if value < 0 || math.IsNaN(value) {
		value = 0
	}

	if value >= 1 || math.IsInf(value, 0) {
		value = 0.99
	}

	return value
}
