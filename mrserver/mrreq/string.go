package mrreq

import (
	"strings"

	"github.com/mondegor/go-webcore/mrcore"
)

const (
	maxLenString = 256
)

// ParseStr - возвращает строковое значение из внешнего строкового параметра по указанному ключу.
// Если параметр пустой, то в зависимости от required возвращается пустая строка или ошибка.
func ParseStr(getter valueGetter, key string, required bool) (string, error) {
	value := strings.TrimSpace(getter.Get(key))

	if value == "" {
		if required {
			return "", mrcore.ErrHttpRequestParamEmpty.New(key)
		}

		return "", nil
	}

	if len(value) > maxLenString {
		return "", mrcore.ErrHttpRequestParamLenMax.New(key, maxLenString)
	}

	return value, nil
}
