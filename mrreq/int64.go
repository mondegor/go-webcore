package mrreq

import (
    "net/http"
    "strconv"

    "github.com/mondegor/go-webcore/mrcore"
)

const (
    maxInt64Len = 32
)

func Int64(r *http.Request, key string) (int64, error) {
    value := r.URL.Query().Get(key)

    if value == "" {
        return 0, mrcore.FactoryErrHttpRequestParamEmpty.New(key)
    }

    if len(value) > maxInt64Len {
        return 0, mrcore.FactoryErrHttpRequestParamLen.New(key, maxInt64Len)
    }

    item, err := strconv.ParseInt(value, 10, 64)

    if err != nil {
        return 0, mrcore.FactoryErrHttpRequestParseParam.New("int64", key, value)
    }

    return item, nil
}
