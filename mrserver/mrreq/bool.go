package mrreq

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/mondegor/go-webcore/mrcore"
)

// ParseRequiredBool - comment func.
func ParseRequiredBool(r *http.Request, key string) (bool, error) {
	value := strings.TrimSpace(r.URL.Query().Get(key))

	if value == "" {
		return false, mrcore.ErrHttpRequestParamEmpty.New(key)
	}

	item, err := strconv.ParseBool(value)
	if err != nil {
		return false, mrcore.ErrHttpRequestParseParam.Wrap(err, key, "RequiredBool", value)
	}

	return item, nil
}

// ParseNullableBool - comment func.
func ParseNullableBool(r *http.Request, key string) (*bool, error) {
	value := strings.TrimSpace(r.URL.Query().Get(key))

	if value == "" {
		return nil, nil //nolint:nilnil
	}

	item, err := strconv.ParseBool(value)
	if err != nil {
		return nil, mrcore.ErrHttpRequestParseParam.Wrap(err, key, "NullableBool", value)
	}

	return &item, nil
}
