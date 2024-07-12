package mrlib

import (
	"math"
)

// RoundFloat - возвращает округлённую версию числа x с указанной точностью.
// Особые случаи:
// - Round(±0) = ±0;
// - Round(±Inf) = ±Inf;
// - Round(NaN) = NaN.
func RoundFloat(x float64, precision int) float64 {
	const baseDecimal = 10
	pow := math.Pow(baseDecimal, float64(precision))

	return math.Round(x*pow) / pow
}

// RoundFloat2 - возвращает RoundFloat с точностью 2.
func RoundFloat2(x float64) float64 {
	return RoundFloat(x, 2)
}

// RoundFloat4 - возвращает RoundFloat с точностью 4.
func RoundFloat4(x float64) float64 {
	return RoundFloat(x, 4)
}

// RoundFloat8 - возвращает RoundFloat с точностью 8.
func RoundFloat8(x float64) float64 {
	return RoundFloat(x, 8)
}
