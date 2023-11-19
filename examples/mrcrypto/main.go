package main

import (
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrcrypto"
)

func main() {
	logger := mrcore.DefaultLogger()

	logger.Info("GenDigitCode: %s", mrcrypto.GenDigitCode(16))
	logger.Info("GenTokenBase64: %s", mrcrypto.GenTokenBase64(16))
	logger.Info("GenTokenHex: %s", mrcrypto.GenTokenHex(16))
	logger.Info("GenDigitCode-with-error: %s", mrcrypto.GenDigitCode(20))

	logger.Info("GenPassword: %s", mrcrypto.GenPassword(64, mrcrypto.PassAll))

	pw := mrcrypto.GenPassword(12, mrcrypto.PassAbc)
	logger.Info("PasswordStrength 12 abc: %s - %s", pw, mrcrypto.PasswordStrength(pw))

	pw = mrcrypto.GenPassword(8, mrcrypto.PassAbcNumerals)
	logger.Info("PasswordStrength 8 abc+num: %s - %s", pw, mrcrypto.PasswordStrength(pw))

	pw = mrcrypto.GenPassword(12, mrcrypto.PassAll)
	logger.Info("PasswordStrength 12 all: %s - %s", pw, mrcrypto.PasswordStrength(pw))
}
