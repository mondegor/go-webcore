package mrcrypt

import (
	"crypto/rand"
	"math/big"
	"strconv"

	"github.com/mondegor/go-webcore/mrlog"
)

func GenDigitCode(length int) string {
	if length < 1 {
		mrlog.Default().Warn().Caller().Msgf("param 'length': %d < 1", length)
		length = 1
	}

	if length > 19 {
		mrlog.Default().Warn().Caller().Msgf("param 'length': %d > 19", length)
		length = 19
	}

	minValue := pow(10, length-1)
	value, err := rand.Int(rand.Reader, big.NewInt(minValue*9))

	if err != nil {
		mrlog.Default().Error().Err(err).Send()
		value = big.NewInt(minValue * 9)
	}

	return strconv.FormatInt(value.Int64()+minValue, 10)
}

func pow(num int64, exponent int) int64 {
	if exponent == 0 {
		return 1
	}

	t := pow(num, exponent/2)

	if exponent%2 == 0 {
		return t * t
	}

	return t * t * num
}
