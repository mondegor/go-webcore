package mrenv

import (
    "net/http"
    "strconv"
    "strings"

    "github.com/mondegor/go-webcore/mrcore"
)

const (
	maxInt64ListLen = 256
)

func Int64ListFromRequest(r *http.Request, key string) ([]int64, error) {
    return parseInt64List(r.URL.Query().Get(key), key)
}

func parseInt64List(name string, value string) ([]int64, error) {
    if value == "" {
        return []int64{}, nil
    }

    if len(value) > maxInt64ListLen {
        return nil, mrcore.FactoryErrHttpRequestParamLen.New(name, maxEnumListLen)
    }

    var items []int64

    for _, item := range strings.Split(value, ",") {
        item = strings.TrimSpace(item)

        i, err := strconv.ParseInt(item, 10, 64)

        if err != nil {
            return nil, mrcore.FactoryErrHttpRequestParseParam.New("int64", name, value)
        }

        items = append(items, i)
    }

    return items, nil
}
