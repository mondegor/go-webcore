package parser

import (
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/mondegor/go-core/mrlog"
	"github.com/mondegor/go-core/mrtype/errors"
)

const (
	typeUUID   = "UUID"
	maxLenUUID = 64
)

type (
	// UUID - парсер UUID.
	UUID struct {
		pathFunc func(r *http.Request, key string) string
		logger   mrlog.Logger
	}
)

// NewUUID - создаёт объект UUID.
func NewUUID(
	pathFunc func(r *http.Request, key string) string,
	logger mrlog.Logger,
) *UUID {
	return &UUID{
		pathFunc: pathFunc,
		logger:   logger,
	}
}

// PathParamUUID - возвращает именованный UUID содержащийся в URL пути.
func (p *UUID) PathParamUUID(r *http.Request, name string) uuid.UUID {
	value, err := uuid.Parse(p.pathFunc(r, name))
	if err != nil {
		p.logger.Warn(r.Context(), "PathParamUUID", "name", name, "error", err)

		return uuid.Nil
	}

	return value
}

// FilterUUID - возвращает UUID поступивший из внешнего запроса.
// Если ключ key не найден или возникнет ошибка, то возвращается uuid.Nil.
func (p *UUID) FilterUUID(r *http.Request, key string) uuid.UUID {
	value, err := p.parseUUID(r.URL.Query().Get(key), false)
	if err != nil {
		p.logger.Warn(r.Context(), "FilterUUID", "key", key, "error", err)

		return uuid.Nil
	}

	return value
}

func (p *UUID) parseUUID(value string, required bool) (uuid.UUID, error) {
	value = strings.TrimSpace(value)

	if value == "" {
		if required {
			return uuid.Nil, errors.NewParamEmptyError(typeUUID)
		}

		return uuid.Nil, nil
	}

	if len(value) > maxLenUUID {
		return uuid.Nil, errors.NewParamLenMaxError(typeUUID, maxLenUUID)
	}

	id, err := uuid.Parse(value)
	if err != nil {
		return uuid.Nil, errors.NewParamIncorrectError(typeUUID, err)
	}

	return id, nil
}
