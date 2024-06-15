package mrreq

import (
	"net/http"
	"strings"

	"github.com/mondegor/go-webcore/mrcore"
)

const (
	maxLenStringsList = 2048
)

// ParseStrList  - comment func.
func ParseStrList(r *http.Request, key string) ([]string, error) {
	value := strings.TrimSpace(r.URL.Query().Get(key))

	if value == "" {
		return nil, nil
	}

	if len(value) > maxLenStringsList {
		return nil, mrcore.ErrHttpRequestParamLenMax.New(key, maxLenStringsList)
	}

	items := strings.Split(value, ",")

	for i, item := range items {
		items[i] = strings.TrimSpace(item)
	}

	return items, nil
}
