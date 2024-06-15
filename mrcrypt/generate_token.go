package mrcrypt

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/mondegor/go-webcore/mrlog"
)

// GenTokenBase64  - comment func.
func GenTokenBase64(length int) string {
	return cutString(base64.StdEncoding.EncodeToString(genToken(length)), length)
}

// GenTokenHex  - comment func.
func GenTokenHex(length int) string {
	return cutString(hex.EncodeToString(genToken(length)), length)
}

// GenTokenHexWithDelimiter  - comment func.
func GenTokenHexWithDelimiter(length, repeat int) string {
	if repeat < 1 {
		mrlog.Default().Warn().Err(fmt.Errorf("param 'repeat': %d < 1", repeat)).Send()
		repeat = 1
	}

	if repeat > 16 {
		mrlog.Default().Warn().Err(fmt.Errorf("param 'repeat': %d > 16", repeat)).Send()
		repeat = 16
	}

	s := make([]string, repeat)

	for i := 0; i < repeat; i++ {
		s[i] = cutString(hex.EncodeToString(genToken(length)), length)
	}

	return strings.Join(s, "-")
}

func genToken(length int) []byte {
	if length < 1 {
		mrlog.Default().Warn().Err(fmt.Errorf("param 'length': %d < 1", length)).Send()
		length = 1
	}

	if length > 256 {
		mrlog.Default().Warn().Err(fmt.Errorf("param 'length': %d > 256", length)).Send()
		length = 256
	}

	value := make([]byte, length)

	if _, err := rand.Read(value); err != nil {
		mrlog.Default().Error().Err(err).Send()

		return nil
	}

	return value
}

func cutString(str string, length int) string {
	if len(str) > length {
		return str[0:length]
	}

	return str
}
