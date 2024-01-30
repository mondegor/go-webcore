package main

import (
	"github.com/mondegor/go-webcore/mrcrypt"
	"github.com/mondegor/go-webcore/mrlog"
)

func main() {
	logger := mrlog.New(mrlog.DebugLevel).With().Str("example", "mrcrypt").Logger()

	logger.Info().Msgf("GenDigitCode: %s", mrcrypt.GenDigitCode(16))
	logger.Info().Msgf("GenTokenBase64: %s", mrcrypt.GenTokenBase64(16))
	logger.Info().Msgf("GenTokenHex: %s", mrcrypt.GenTokenHex(16))
	logger.Info().Msgf("GenDigitCode-with-error: %s", mrcrypt.GenDigitCode(20))

	logger.Info().Msgf("GenPassword: %s", mrcrypt.GenPassword(64, mrcrypt.PassAll))

	pw := mrcrypt.GenPassword(12, mrcrypt.PassAbc)
	logger.Info().Msgf("PasswordStrength 12 abc: %s - %s", pw, mrcrypt.PasswordStrength(pw))

	pw = mrcrypt.GenPassword(8, mrcrypt.PassAbcNumerals)
	logger.Info().Msgf("PasswordStrength 8 abc+num: %s - %s", pw, mrcrypt.PasswordStrength(pw))

	pw = mrcrypt.GenPassword(12, mrcrypt.PassAll)
	logger.Info().Msgf("PasswordStrength 12 all: %s - %s", pw, mrcrypt.PasswordStrength(pw))
}
