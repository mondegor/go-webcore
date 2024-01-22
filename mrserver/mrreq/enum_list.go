package mrreq

import (
	"net/http"
	"strings"

	"github.com/mondegor/go-webcore/mrcore"
)

const (
	maxLenEnumList = 256
)

func ParseEnumList(r *http.Request, key string) ([]string, error) {
	value := strings.TrimSpace(r.URL.Query().Get(key))

	if value == "" {
		return []string{}, nil
	}

	if len(value) > maxLenEnumList {
		return nil, mrcore.FactoryErrHttpRequestParamLenMax.New(key, maxLenEnumList)
	}

	items := strings.Split(strings.ToUpper(value), ",")

	for i, item := range items {
		item = strings.TrimSpace(item)

		if !regexpEnum.MatchString(item) {
			return nil, mrcore.FactoryErrHttpRequestParseParam.New(key, "Enum", value)
		}

		items[i] = item
	}

	return items, nil
}
