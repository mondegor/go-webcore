package main

import (
	"github.com/mondegor/go-webcore/mrcrypt"
	"github.com/mondegor/go-webcore/mrlog"
)

func main() {
	logger := mrlog.Default().With().Str("example", "mrcrypt").Logger()

	lib := mrcrypt.New(logger)

	logger.Info().Msgf("GenDigitCode: %s", lib.GenDigitCode(16))
	logger.Info().Msgf("GenTokenBase64: %s", lib.GenTokenBase64(16))
	logger.Info().Msgf("GenTokenHex: %s", lib.GenTokenHex(16))
	logger.Info().Msgf("GenDigitCode-with-error: %s", lib.GenDigitCode(20))

	logger.Info().Msgf("GenPassword: %s", lib.GenPassword(64, mrcrypt.PassAll))

	pw := lib.GenPassword(12, mrcrypt.PassAbc)
	logger.Info().Msgf("PasswordStrength 12 abc: %s - %s", pw, lib.PasswordStrength(pw))

	pw = lib.GenPassword(8, mrcrypt.PassAbcNumerals)
	logger.Info().Msgf("PasswordStrength 8 abc+num: %s - %s", pw, lib.PasswordStrength(pw))

	pw = lib.GenPassword(12, mrcrypt.PassAll)
	logger.Info().Msgf("PasswordStrength 12 all: %s - %s", pw, lib.PasswordStrength(pw))
}
