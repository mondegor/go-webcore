package mrreq

import (
    "net/http"
    "strconv"
    "strings"

    "github.com/mondegor/go-webcore/mrcore"
)

const (
	maxInt64ListLen = 256
)

func Int64List(r *http.Request, key string) ([]int64, error) {
    value := r.URL.Query().Get(key)

    if value == "" {
        return []int64{}, nil
    }

    if len(value) > maxInt64ListLen {
        return nil, mrcore.FactoryErrHttpRequestParamLen.New(key, maxEnumListLen)
    }

    var items []int64

    for _, item := range strings.Split(value, ",") {
        item = strings.TrimSpace(item)

        i, err := strconv.ParseInt(item, 10, 64)

        if err != nil {
            return nil, mrcore.FactoryErrHttpRequestParseParam.New("int64", key, value)
        }

        items = append(items, i)
    }

    return items, nil
}
