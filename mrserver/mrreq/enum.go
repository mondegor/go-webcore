package mrreq

import (
	"regexp"
	"strings"

	"github.com/mondegor/go-sysmess/mrerr/mr"
)

const (
	maxLenEnum = 64
)

var regexpEnum = regexp.MustCompile(`^[A-Z]([A-Z0-9_]+)?[A-Z0-9]$`)

// ParseEnum - возвращает Enum значение из внешнего строкового параметра по указанному ключу.
// Если параметр пустой, то в зависимости от required возвращается пустая строка или ошибка.
func ParseEnum(getter valueGetter, key string, required bool) (string, error) {
	value := strings.TrimSpace(getter.Get(key))

	if value == "" {
		if required {
			return "", mr.ErrHttpRequestParamEmpty.New(key)
		}

		return "", nil
	}

	if len(value) > maxLenEnum {
		return "", mr.ErrHttpRequestParamLenMax.New(key, maxLenEnum)
	}

	value = strings.ToUpper(value)

	if !regexpEnum.MatchString(value) {
		return "", mr.ErrHttpRequestParseParam.New(key, "Enum", value)
	}

	return value, nil
}
