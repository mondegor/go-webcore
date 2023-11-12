package mrreq

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/mondegor/go-webcore/mrcore"
)

const (
	maxLenInt64 = 32
)

func ParseInt64(r *http.Request, key string, required bool) (int64, error) {
	value := strings.TrimSpace(r.URL.Query().Get(key))

	if value == "" {
		if required {
			return 0, mrcore.FactoryErrHttpRequestParamEmpty.New(key)
		}

		return 0, nil
	}

	if len(value) > maxLenInt64 {
		return 0, mrcore.FactoryErrHttpRequestParamLenMax.New(key, maxLenInt64)
	}

	item, err := strconv.ParseInt(value, 10, 64)

	if err != nil {
		return 0, mrcore.FactoryErrHttpRequestParseParam.New("Int64", key, value)
	}

	return item, nil
}
