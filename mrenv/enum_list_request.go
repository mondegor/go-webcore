package mrenv

import (
    "net/http"
    "strings"

    "github.com/mondegor/go-webcore/mrcore"
)

const (
	maxEnumListLen = 256
)

func EnumListFromRequest(r *http.Request, key string) ([]string, error) {
    return parseEnumList(r.URL.Query().Get(key), key)
}

func parseEnumList(name string, value string) ([]string, error) {
    if value == "" {
        return []string{}, nil
    }

    if len(value) > maxEnumListLen {
        return nil, mrcore.FactoryErrHttpRequestParamLen.New(name, maxEnumListLen)
    }

    var items []string

    for _, item := range strings.Split(strings.ToUpper(value), ",") {
        item = strings.TrimSpace(item)

        if !regexpEnum.MatchString(item) {
            return nil, mrcore.FactoryErrHttpRequestParseParam.New("enum", name, value)
        }

        items = append(items, item)
    }

    return items, nil
}
