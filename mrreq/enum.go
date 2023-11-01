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

func ParseEnum(r *http.Request, key string, required bool) (string, error) {
    value := r.URL.Query().Get(key)

    if value == "" {
        if required {
            return "", mrcore.FactoryErrHttpRequestParamEmpty.New(key)
        }

        return "", nil
    }

    if len(value) > maxEnumLen {
        return "", mrcore.FactoryErrHttpRequestParamLenMax.New(key, maxEnumLen)
    }

    value = strings.ToUpper(strings.TrimSpace(value))

    if !regexpEnum.MatchString(value) {
        return "", mrcore.FactoryErrHttpRequestParseParam.New("Enum", key, value)
    }

    return value, nil
}
