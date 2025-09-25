package mrreq

import (
	"strconv"
	"strings"

	"github.com/mondegor/go-sysmess/mrerr/mr"
)

const (
	maxLenUint64 = 32
)

// ParseUint64 - возвращает Uint64 значение из внешнего строкового параметра по указанному ключу.
// Если параметр пустой, то в зависимости от required возвращается 0 или ошибка.
func ParseUint64(getter valueGetter, key string, required bool) (uint64, error) {
	value := strings.TrimSpace(getter.Get(key))

	if value == "" {
		if required {
			return 0, mr.ErrHttpRequestParamEmpty.New(key)
		}

		return 0, nil
	}

	if len(value) > maxLenUint64 {
		return 0, mr.ErrHttpRequestParamLenMax.New(key, maxLenUint64)
	}

	item, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return 0, mr.ErrHttpRequestParseParam.Wrap(err, key, "Uint64", value)
	}

	return item, nil
}
