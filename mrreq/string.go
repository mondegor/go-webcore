package mrreq

import (
    "net/http"
    "strings"

    "github.com/mondegor/go-webcore/mrcore"
)

const (
    maxStringLen = 128
)

func ParseStr(r *http.Request, key string, required bool) (string, error) {
    value := r.URL.Query().Get(key)

    if value == "" {
        if required {
            return "", mrcore.FactoryErrHttpRequestParamEmpty.New(key)
        }

        return "", nil
    }

    if len(value) > maxStringLen {
        return "", mrcore.FactoryErrHttpRequestParamLenMax.New(key, maxStringLen)
    }

    return strings.ToUpper(strings.TrimSpace(value)), nil
}
