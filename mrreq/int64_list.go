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

func ParseInt64List(r *http.Request, key string) ([]int64, error) {
    value := r.URL.Query().Get(key)

    if value == "" {
        return []int64{}, nil
    }

    if len(value) > maxInt64ListLen {
        return nil, mrcore.FactoryErrHttpRequestParamLenMax.New(key, maxEnumListLen)
    }

    itemsTmp := strings.Split(value, ",")
    items := make([]int64, len(itemsTmp))

    for i, item := range itemsTmp {
        itemN, err := strconv.ParseInt(strings.TrimSpace(item), 10, 64)

        if err != nil {
            return nil, mrcore.FactoryErrHttpRequestParseParam.New("Int64", key, value)
        }

        items[i] = itemN
    }

    return items, nil
}
