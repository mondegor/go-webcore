package mrreq

import (
    "net/http"
    "regexp"
    "strings"

    "github.com/mondegor/go-webcore/mrcore"
)

const (
    maxEnumLen = 64
)

var (
    regexpEnum = regexp.MustCompile(`^[A-Z]([A-Z0-9_]+)?[A-Z0-9]$`)
)

func Enum(r *http.Request, key string) (string, error) {
    value := r.URL.Query().Get(key)

    if value == "" {
        return "", mrcore.FactoryErrHttpRequestParamEmpty.New(key)
    }

    if len(value) > maxEnumLen {
        return "", mrcore.FactoryErrHttpRequestParamLen.New(key, maxEnumLen)
    }

    value = strings.ToUpper(strings.TrimSpace(value))

    if !regexpEnum.MatchString(value) {
        return "", mrcore.FactoryErrHttpRequestParseParam.New("enum", key, value)
    }

    return value, nil
}
