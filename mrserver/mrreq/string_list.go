package mrreq

import (
	"strings"

	"github.com/mondegor/go-sysmess/mrerr/mr"
)

const (
	maxLenStringsList = 2048
)

// ParseStrList - возвращает массив строковых значений из внешнего строкового параметра по указанному ключу.
// Если параметр пустой, то возвращается пустой массив.
func ParseStrList(getter valueGetter, key string) ([]string, error) {
	value := strings.TrimSpace(getter.Get(key))

	if value == "" {
		return nil, nil
	}

	if len(value) > maxLenStringsList {
		return nil, mr.ErrHttpRequestParamLenMax.New(key, maxLenStringsList)
	}

	items := strings.Split(value, ",")

	for i, item := range items {
		items[i] = strings.TrimSpace(item)
	}

	return items, nil
}
