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
}
