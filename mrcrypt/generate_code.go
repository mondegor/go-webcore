package mrcrypt

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strconv"

	"github.com/mondegor/go-webcore/mrlog"
)

// GenDigitCode  - comment func.
func GenDigitCode(length int) string {
	if length < 1 {
		mrlog.Default().Warn().Err(fmt.Errorf("param 'length': %d < 1", length)).Send()
		length = 1
	}

	if length > 19 {
		mrlog.Default().Warn().Err(fmt.Errorf("param 'length': %d > 19", length)).Send()
		length = 19
	}

	minValue := pow(10, length-1)

	value, err := rand.Int(rand.Reader, big.NewInt(minValue*9))
	if err != nil {
		value = big.NewInt(minValue * 9)

		mrlog.Default().Error().Err(err).Send()
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
