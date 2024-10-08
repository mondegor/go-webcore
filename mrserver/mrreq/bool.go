package mrreq

import (
	"strconv"
	"strings"

	"github.com/mondegor/go-webcore/mrcore"
)

// ParseRequiredBool - возвращает Bool значение из внешнего строкового параметра по указанному ключу.
// Если параметр пустой, то возвращается ошибка.
func ParseRequiredBool(getter valueGetter, key string) (bool, error) {
	value := strings.TrimSpace(getter.Get(key))

	if value == "" {
		return false, mrcore.ErrHttpRequestParamEmpty.New(key)
	}

	item, err := strconv.ParseBool(value)
	if err != nil {
		return false, mrcore.ErrHttpRequestParseParam.Wrap(err, key, "RequiredBool", value)
	}

	return item, nil
}

// ParseNullableBool - возвращает Bool значение из внешнего строкового параметра по указанному ключу.
// Если параметр пустой, то возвращается nil.
func ParseNullableBool(getter valueGetter, key string) (*bool, error) {
	value := strings.TrimSpace(getter.Get(key))

	if value == "" {
		return nil, nil //nolint:nilnil
	}

	item, err := strconv.ParseBool(value)
	if err != nil {
		return nil, mrcore.ErrHttpRequestParseParam.Wrap(err, key, "NullableBool", value)
	}

	return &item, nil
}
