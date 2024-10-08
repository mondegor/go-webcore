package mrreq

import (
	"strings"

	"github.com/google/uuid"

	"github.com/mondegor/go-webcore/mrcore"
)

const (
	maxLenUUID = 64
)

// ParseUUID - возвращает UUID значение из внешнего строкового параметра по указанному ключу.
// Если параметр пустой, то в зависимости от required возвращается нулевой UUID или ошибка.
func ParseUUID(getter valueGetter, key string, required bool) (uuid.UUID, error) {
	value := strings.TrimSpace(getter.Get(key))

	if value == "" {
		if required {
			return uuid.Nil, mrcore.ErrHttpRequestParamEmpty.New(key)
		}

		return uuid.Nil, nil
	}

	if len(value) > maxLenUUID {
		return uuid.Nil, mrcore.ErrHttpRequestParamLenMax.New(key, maxLenInt64)
	}

	item, err := uuid.Parse(value)
	if err != nil {
		return uuid.Nil, mrcore.ErrHttpRequestParseParam.Wrap(err, key, "UUID", value)
	}

	return item, nil
}
