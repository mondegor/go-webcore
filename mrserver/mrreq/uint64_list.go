package mrreq

import (
	"strconv"
	"strings"

	"github.com/mondegor/go-webcore/mrcore"
)

const (
	maxLenUint64List = 256
)

// ParseUint64List - возвращает массив Uint64 значений из внешнего строкового параметра по указанному ключу.
// Если параметр пустой, то возвращается пустой массив.
func ParseUint64List(getter valueGetter, key string) ([]uint64, error) {
	value := strings.TrimSpace(getter.Get(key))

	if value == "" {
		return nil, nil
	}

	if len(value) > maxLenUint64List {
		return nil, mrcore.ErrHttpRequestParamLenMax.New(key, maxLenUint64List)
	}

	itemsTmp := strings.Split(value, ",")
	items := make([]uint64, 0, len(itemsTmp))

	for i := range itemsTmp {
		item, err := strconv.ParseUint(strings.TrimSpace(itemsTmp[i]), 10, 64)
		if err != nil {
			return nil, mrcore.ErrHttpRequestParseParam.Wrap(err, key, "Uint64", value)
		}

		items = append(items, item)
	}

	return items, nil
}
