package mrcrypto

import (
    "crypto/rand"
    "encoding/base64"
    "encoding/hex"
    "log"
    "strings"

    "github.com/mondegor/go-webcore/mrcore"
)

func GenTokenBase64(length int) string {
    return base64.StdEncoding.EncodeToString(GenToken(length))
}

func GenTokenHex(length int) string {
    return hex.EncodeToString(GenToken(length))
}

func GenTokenHexWithDelimiter(length int, repeat int) string {
    if repeat < 1 {
        mrcore.LogError("param 'repeat': %d < 1", repeat)
        repeat = 1
    }

    if repeat > 16 {
        mrcore.LogError("param 'repeat': %d > 16", repeat)
        repeat = 16
    }

    s := make([]string, repeat)

    for i := 0; i < repeat; i++ {
        s[i] = hex.EncodeToString(GenToken(length))
    }

    return strings.Join(s, "-")
}

func GenToken(length int) []byte {
    if length < 1 {
        mrcore.LogError("param 'length': %d < 1", length)
        length = 1
    }

    if length > 256 {
        mrcore.LogError("param 'length': %d > 256", length)
        length = 256
    }

    value := make([]byte, length)

    _, err := rand.Read(value)

    if err != nil {
        log.Print(err)
        return []byte{}
    }

    return value
}
