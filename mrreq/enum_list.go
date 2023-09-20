package mrreq

import (
    "net/http"
    "strings"

    "github.com/mondegor/go-webcore/mrcore"
)

const (
    maxEnumListLen = 256
)

func EnumList(r *http.Request, key string) ([]string, error) {
    value := r.URL.Query().Get(key)

    if value == "" {
        return []string{}, nil
    }

    if len(value) > maxEnumListLen {
        return nil, mrcore.FactoryErrHttpRequestParamLen.New(key, maxEnumListLen)
    }

    var items []string

    for _, item := range strings.Split(strings.ToUpper(value), ",") {
        item = strings.TrimSpace(item)

        if !regexpEnum.MatchString(item) {
            return nil, mrcore.FactoryErrHttpRequestParseParam.New("enum", key, value)
        }

        items = append(items, item)
    }

    return items, nil
}
