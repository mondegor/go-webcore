package mrreq

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/mondegor/go-webcore/mrcore"
)

func ParseRequiredBool(r *http.Request, key string) (bool, error) {
	value := strings.TrimSpace(r.URL.Query().Get(key))

	if value == "" {
		return false, mrcore.FactoryErrHttpRequestParamEmpty.New(key)
	}

	item, err := strconv.ParseBool(value)

	if err != nil {
		return false, mrcore.FactoryErrHttpRequestParseParam.New(key, "RequiredBool", value)
	}

	return item, nil
}

func ParseNullableBool(r *http.Request, key string) (*bool, error) {
	value := strings.TrimSpace(r.URL.Query().Get(key))

	if value == "" {
		return nil, nil
	}

	item, err := strconv.ParseBool(value)

	if err != nil {
		return nil, mrcore.FactoryErrHttpRequestParseParam.New(key, "NullableBool", value)
	}

	return &item, nil
}
