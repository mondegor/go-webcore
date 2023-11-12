package mrreq

import (
	"net/http"
	"strings"

	"github.com/mondegor/go-webcore/mrcore"
)

const (
	maxLenStringsList = 1024
)

func ParseStrList(r *http.Request, key string) ([]string, error) {
	value := strings.TrimSpace(r.URL.Query().Get(key))

	if value == "" {
		return []string{}, nil
	}

	if len(value) > maxLenStringsList {
		return nil, mrcore.FactoryErrHttpRequestParamLenMax.New(key, maxLenStringsList)
	}

	items := strings.Split(value, ",")

	for i, item := range items {
		items[i] = strings.TrimSpace(item)
	}

	return items, nil
}
