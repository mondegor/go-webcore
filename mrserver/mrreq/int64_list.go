package mrreq

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/mondegor/go-webcore/mrcore"
)

const (
	maxLenInt64List = 256
)

// ParseInt64List  - comment func.
func ParseInt64List(r *http.Request, key string) ([]int64, error) {
	value := strings.TrimSpace(r.URL.Query().Get(key))

	if value == "" {
		return nil, nil
	}

	if len(value) > maxLenInt64List {
		return nil, mrcore.ErrHttpRequestParamLenMax.New(key, maxLenInt64List)
	}

	itemsTmp := strings.Split(value, ",")
	items := make([]int64, len(itemsTmp))

	for i, item := range itemsTmp {
		itemN, err := strconv.ParseInt(strings.TrimSpace(item), 10, 64)
		if err != nil {
			return nil, mrcore.ErrHttpRequestParseParam.Wrap(err, key, "Int64", value)
		}

		items[i] = itemN
	}

	return items, nil
}
