package mrreq

import (
	"strconv"
	"strings"

	"github.com/mondegor/go-sysmess/mrerr/mr"
)

const (
	maxLenFloat64 = 64
)

// ParseFloat64 - возвращает Float64 значение из внешнего строкового параметра по указанному ключу.
// Если параметр пустой, то в зависимости от required возвращается 0 или ошибка.
func ParseFloat64(getter valueGetter, key string, required bool) (float64, error) {
	value := strings.TrimSpace(getter.Get(key))

	if value == "" {
		if required {
			return 0, mr.ErrHttpRequestParamEmpty.New(key)
		}

		return 0, nil
	}

	if len(value) > maxLenFloat64 {
		return 0, mr.ErrHttpRequestParamLenMax.New(key, maxLenFloat64)
	}

	item, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, mr.ErrHttpRequestParseParam.Wrap(err, key, "Float64", value)
	}

	return item, nil
}
