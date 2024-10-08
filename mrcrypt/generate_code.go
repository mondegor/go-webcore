package mrcrypt

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strconv"
)

// GenDigitCode - comment func.
func (l *Lib) GenDigitCode(length int) string {
	if length < 1 {
		l.logger.Warn().Err(fmt.Errorf("param 'length': %d < 1", length)).Send()
		length = 1
	}

	if length > 19 {
		l.logger.Warn().Err(fmt.Errorf("param 'length': %d > 19", length)).Send()
		length = 19
	}

	minValue := l.pow(10, length-1)

	value, err := rand.Int(rand.Reader, big.NewInt(minValue*9))
	if err != nil {
		value = big.NewInt(minValue * 9)

		l.logger.Error().Err(err).Send()
	}

	return strconv.FormatInt(value.Int64()+minValue, 10)
}

func (l *Lib) pow(num int64, exponent int) int64 {
	if exponent < 1 {
		return 1
	}

	t := l.pow(num, exponent/2)

	if exponent%2 == 0 {
		return t * t
	}

	return t * t * num
}
