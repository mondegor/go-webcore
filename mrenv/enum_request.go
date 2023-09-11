package mrenv

import (
    "net/http"
    "regexp"
    "strings"

    "github.com/mondegor/go-webcore/mrcore"
)

const (
	maxEnumLen = 32
)

var (
	regexpEnum = regexp.MustCompile(`^[A-Z]([A-Z0-9_]+)?[A-Z0-9]$`)
)

func EnumFromRequest(r *http.Request, key string) (string, error) {
    return parseEnum(key, r.URL.Query().Get(key))
}

func parseEnum(name string, value string) (string, error) {
    if value == "" {
        return "", nil
    }

    if len(value) > maxEnumLen {
        return "", mrcore.FactoryErrHttpRequestParamLen.New(name, maxEnumLen)
    }

    value = strings.ToUpper(strings.TrimSpace(value))

    if !regexpEnum.MatchString(value) {
        return "", mrcore.FactoryErrHttpRequestParseParam.New("enum", name, value)
    }

    return value, nil
}
