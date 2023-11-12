package mrreq

import (
    "net/http"
    "regexp"
    "strings"

    "github.com/mondegor/go-webcore/mrcore"
)

const (
    maxLenEnum = 64
)

var (
    regexpEnum = regexp.MustCompile(`^[A-Z]([A-Z0-9_]+)?[A-Z0-9]$`)
)

func ParseEnum(r *http.Request, key string, required bool) (string, error) {
    value := strings.TrimSpace(r.URL.Query().Get(key))

    if value == "" {
        if required {
            return "", mrcore.FactoryErrHttpRequestParamEmpty.New(key)
        }

        return "", nil
    }

    if len(value) > maxLenEnum {
        return "", mrcore.FactoryErrHttpRequestParamLenMax.New(key, maxLenEnum)
    }

    value = strings.ToUpper(value)

    if !regexpEnum.MatchString(value) {
        return "", mrcore.FactoryErrHttpRequestParseParam.New("Enum", key, value)
    }

    return value, nil
}
