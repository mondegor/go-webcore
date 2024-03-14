package mrcrypt

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"log"
	"strings"

	"github.com/mondegor/go-webcore/mrlog"
)

func GenTokenBase64(length int) string {
	return cutString(base64.StdEncoding.EncodeToString(genToken(length)), length)
}

func GenTokenHex(length int) string {
	return cutString(hex.EncodeToString(genToken(length)), length)
}

func GenTokenHexWithDelimiter(length, repeat int) string {
	if repeat < 1 {
		mrlog.Default().Warn().Caller().Msgf("param 'repeat': %d < 1", repeat)
		repeat = 1
	}

	if repeat > 16 {
		mrlog.Default().Warn().Caller().Msgf("param 'repeat': %d > 16", repeat)
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
		mrlog.Default().Warn().Caller(1).Msgf("param 'length': %d < 1", length)
		length = 1
	}

	if length > 256 {
		mrlog.Default().Warn().Caller(1).Msgf("param 'length': %d > 256", length)
		length = 256
	}

	value := make([]byte, length)

	if _, err := rand.Read(value); err != nil {
		log.Print(err)
		return []byte{}
	}

	return value
}

func cutString(str string, length int) string {
	if len(str) > length {
		return str[0:length]
	}

	return str
}
