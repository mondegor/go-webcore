package mrreq

import (
	"strconv"
	"strings"

	"github.com/mondegor/go-sysmess/mrerr/mr"
)

const (
	maxLenInt64List = 256
)

// ParseInt64List - возвращает массив Int64 значений из внешнего строкового параметра по указанному ключу.
// Если параметр пустой, то возвращается пустой массив.
func ParseInt64List(getter valueGetter, key string) ([]int64, error) {
	value := strings.TrimSpace(getter.Get(key))

	if value == "" {
		return nil, nil
	}

	if len(value) > maxLenInt64List {
		return nil, mr.ErrHttpRequestParamLenMax.New(key, maxLenInt64List)
	}

	itemsTmp := strings.Split(value, ",")
	items := make([]int64, 0, len(itemsTmp))

	for i := range itemsTmp {
		item, err := strconv.ParseInt(strings.TrimSpace(itemsTmp[i]), 10, 64)
		if err != nil {
			return nil, mr.ErrHttpRequestParseParam.Wrap(err, key, "Int64", value)
		}

		items = append(items, item)
	}

	return items, nil
}
