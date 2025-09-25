package mrreq

import (
	"strings"
	"time"

	"github.com/mondegor/go-sysmess/mrerr/mr"
)

const (
	maxLenDateTime = 64
)

// ParseDateTime - возвращает time.Time значение из внешнего строкового параметра по указанному ключу.
// Значение строкового параметра должно быть указано в формате RFC3339.
// Если параметр пустой, то в зависимости от required возвращается нулевое время или ошибка.
func ParseDateTime(getter valueGetter, key string, required bool) (time.Time, error) {
	value := strings.TrimSpace(getter.Get(key))

	if value == "" {
		if required {
			return time.Time{}, mr.ErrHttpRequestParamEmpty.New(key)
		}

		return time.Time{}, nil
	}

	if len(value) > maxLenDateTime {
		return time.Time{}, mr.ErrHttpRequestParamLenMax.New(key, maxLenDateTime)
	}

	item, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return time.Time{}, mr.ErrHttpRequestParseParam.Wrap(err, key, "DateTime", value)
	}

	return item, nil
}
