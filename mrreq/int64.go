package mrreq

import (
    "net/http"
    "strconv"

    "github.com/mondegor/go-webcore/mrcore"
)

const (
    maxInt64Len = 32
)

func ParseInt64(r *http.Request, key string, required bool) (int64, error) {
    value := r.URL.Query().Get(key)

    if value == "" {
        if required {
            return 0, mrcore.FactoryErrHttpRequestParamEmpty.New(key)
        }

        return 0, nil
    }

    if len(value) > maxInt64Len {
        return 0, mrcore.FactoryErrHttpRequestParamLenMax.New(key, maxInt64Len)
    }

    item, err := strconv.ParseInt(value, 10, 64)

    if err != nil {
        return 0, mrcore.FactoryErrHttpRequestParseParam.New("Int64", key, value)
    }

    return item, nil
}
