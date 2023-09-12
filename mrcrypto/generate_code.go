package mrcrypto

import (
    "crypto/rand"
    "fmt"
    "math/big"

    "github.com/mondegor/go-webcore/mrcore"
)

func GenDigitCode(length int) string {
    if length < 1 {
        mrcore.LogError("param 'length': %d < 1", length)
        length = 1
    }

    if length > 19 {
        mrcore.LogError("param 'length': %d > 19", length)
        length = 19
    }

    min := pow(10, length - 1)
    value, err := rand.Int(rand.Reader, big.NewInt(min * 9))

    if err != nil {
        mrcore.LogErr(err)
        value = big.NewInt(min * 9)
    }

    return fmt.Sprintf("%d", value.Int64() + min)
}

func pow(num int64, exponent int) int64 {
    if exponent == 0 {
        return 1
    }

    t := pow(num, exponent / 2)

    if exponent % 2 == 0 {
        return t * t
    }

    return t * t * num
}
