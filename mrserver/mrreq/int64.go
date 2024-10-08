package mrreq

import (
	"strconv"
	"strings"

	"github.com/mondegor/go-webcore/mrcore"
)

const (
	maxLenInt64 = 32
)

// ParseInt64 - возвращает Int64 значение из внешнего строкового параметра по указанному ключу.
// Если параметр пустой, то в зависимости от required возвращается 0 или ошибка.
func ParseInt64(getter valueGetter, key string, required bool) (int64, error) {
	value := strings.TrimSpace(getter.Get(key))

	if value == "" {
		if required {
			return 0, mrcore.ErrHttpRequestParamEmpty.New(key)
		}

		return 0, nil
	}

	if len(value) > maxLenInt64 {
		return 0, mrcore.ErrHttpRequestParamLenMax.New(key, maxLenInt64)
	}

	item, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, mrcore.ErrHttpRequestParseParam.Wrap(err, "Int64", value)
	}

	return item, nil
}
