package mrcrypto

import (
    "crypto/rand"
    "encoding/base64"
    "encoding/hex"
    "strings"

    "github.com/mondegor/go-core/mrlog"
)

func GenTokenBase64(length int) string {
    return base64.StdEncoding.EncodeToString(GenToken(length))
}

func GenTokenHex(length int) string {
    return hex.EncodeToString(GenToken(length))
}

func GenTokenHexWithDelimiter(length int, repeat int) string {
    if repeat < 1 {
        mrlog.Error("param repeat < 1")
        repeat = 1
    }

    if repeat > 16 {
        mrlog.Error("param repeat > 16")
        repeat = 16
    }

    var s []string

    for i := 0; i < repeat; i++ {
        s = append(s, hex.EncodeToString(GenToken(length)))
    }

    return strings.Join(s, "-")
}

func GenToken(length int) []byte {
    if length < 1 {
        mrlog.Error("param length < 1")
        length = 1
    }

    if length > 256 {
        mrlog.Error("param length > 256")
        length = 256
    }

    value := make([]byte, length)

    _, err := rand.Read(value)

    if err != nil {
        mrlog.Error(err)
        return []byte{}
    }

    return value
}
