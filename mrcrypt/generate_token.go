package mrcrypt

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
)

// GenTokenBase64 - comment func.
func (l *Lib) GenTokenBase64(length int) string {
	return l.cutString(base64.StdEncoding.EncodeToString(l.genToken(length)), length)
}

// GenTokenHex - comment func.
func (l *Lib) GenTokenHex(length int) string {
	return l.cutString(hex.EncodeToString(l.genToken(length)), length)
}

// GenTokenHexWithDelimiter - comment func.
func (l *Lib) GenTokenHexWithDelimiter(length, repeat int) string {
	if repeat < 1 {
		l.logger.Warn().Err(fmt.Errorf("param 'repeat': %d < 1", repeat)).Send()
		repeat = 1
	}

	if repeat > 16 {
		l.logger.Warn().Err(fmt.Errorf("param 'repeat': %d > 16", repeat)).Send()
		repeat = 16
	}

	s := make([]string, repeat)

	for i := 0; i < repeat; i++ {
		s[i] = l.cutString(hex.EncodeToString(l.genToken(length)), length)
	}

	return strings.Join(s, "-")
}

func (l *Lib) genToken(length int) []byte {
	if length < 1 {
		l.logger.Warn().Err(fmt.Errorf("param 'length': %d < 1", length)).Send()
		length = 1
	}

	if length > 256 {
		l.logger.Warn().Err(fmt.Errorf("param 'length': %d > 256", length)).Send()
		length = 256
	}

	value := make([]byte, length)

	if _, err := rand.Read(value); err != nil {
		l.logger.Error().Err(err).Send()

		return nil
	}

	return value
}

func (l *Lib) cutString(str string, length int) string {
	if len(str) > length {
		return str[0:length]
	}

	return str
}
