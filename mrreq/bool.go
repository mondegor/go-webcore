package mrreq

import (
    "net/http"
    "strconv"
    "strings"

    "github.com/mondegor/go-webcore/mrcore"
    "github.com/mondegor/go-webcore/mrtype"
)

func ParseRequiredBool(r *http.Request, key string) (bool, error) {
    value := strings.TrimSpace(r.URL.Query().Get(key))

    if value == "" {
        return false, mrcore.FactoryErrHttpRequestParamEmpty.New(key)
    }

    item, err := strconv.ParseBool(value)

    if err != nil {
        return false, mrcore.FactoryErrHttpRequestParseParam.New("RequiredBool", key, value)
    }

    return item, nil
}

func ParseNullableBool(r *http.Request, key string) (mrtype.NullableBool, error) {
    value := strings.TrimSpace(r.URL.Query().Get(key))

    if value == "" {
        return mrtype.NullableBoolNull, nil
    }

    item, err := strconv.ParseBool(value)

    if err != nil {
        return mrtype.NullableBoolNull, mrcore.FactoryErrHttpRequestParseParam.New("NullableBool", key, value)
    }

    if item {
        return mrtype.NullableBoolTrue, nil
    }

    return mrtype.NullableBoolFalse, nil
}
