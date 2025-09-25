package mrlib

import (
	"math"
)

// Точность float числа (диаметр окрестности).
const (
	EqualityThresholdE3 = 1e-3 // 1 / 1000
	EqualityThresholdE6 = 1e-6 // 1 / 1000000
	EqualityThresholdE9 = 1e-9 // 1 / 1000000000
)

// EqualFloat - сообщает, находятся ли два float числа в одной окрестности указанного диаметра.
func EqualFloat(first, second, equalityThreshold float64) bool {
	return math.Abs(first-second) <= equalityThreshold
}

// EqualFloatE3 - сообщает, находятся ли два float числа в одной окрестности EqualityThresholdE3.
func EqualFloatE3(first, second float64) bool {
	return EqualFloat(first, second, EqualityThresholdE3)
}

// EqualFloatE6 - сообщает, находятся ли два float числа в одной окрестности EqualityThresholdE6.
func EqualFloatE6(first, second float64) bool {
	return EqualFloat(first, second, EqualityThresholdE6)
}

// EqualFloatE9 - сообщает, находятся ли два float числа в одной окрестности EqualityThresholdE9.
func EqualFloatE9(first, second float64) bool {
	return EqualFloat(first, second, EqualityThresholdE9)
}
