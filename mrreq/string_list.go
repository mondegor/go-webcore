package mrreq

import (
    "net/http"
    "strings"

    "github.com/mondegor/go-webcore/mrcore"
)

const (
    maxStringsListLen = 1024
)

func ParseStrList(r *http.Request, key string) ([]string, error) {
    value := r.URL.Query().Get(key)

    if value == "" {
        return []string{}, nil
    }

    if len(value) > maxStringsListLen {
        return nil, mrcore.FactoryErrHttpRequestParamLenMax.New(key, maxStringsListLen)
    }

    items := strings.Split(strings.ToUpper(value), ",")

    for i, item := range items {
        items[i] = strings.TrimSpace(item)
    }

    return items, nil
}
