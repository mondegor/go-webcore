package mrreq

import (
	"strings"

	"github.com/google/uuid"
	"github.com/mondegor/go-sysmess/mrerr/mr"
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
			return uuid.Nil, mr.ErrHttpRequestParamEmpty.New(key)
		}

		return uuid.Nil, nil
	}

	if len(value) > maxLenUUID {
		return uuid.Nil, mr.ErrHttpRequestParamLenMax.New(key, maxLenInt64)
	}

	item, err := uuid.Parse(value)
	if err != nil {
		return uuid.Nil, mr.ErrHttpRequestParseParam.Wrap(err, key, "UUID", value)
	}

	return item, nil
}
