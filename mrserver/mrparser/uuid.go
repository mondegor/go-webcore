package mrparser

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/mondegor/go-sysmess/mrlog"

	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/mrreq"
)

type (
	// UUID - парсер UUID.
	UUID struct {
		pathFunc mrserver.RequestParserParamFunc
		logger   mrlog.Logger
	}
)

// NewUUID - создаёт объект UUID.
func NewUUID(pathFunc mrserver.RequestParserParamFunc, logger mrlog.Logger) *UUID {
	return &UUID{
		pathFunc: pathFunc,
		logger:   logger,
	}
}

// PathParamUUID - возвращает именованный UUID содержащийся в URL пути.
func (p *UUID) PathParamUUID(r *http.Request, name string) uuid.UUID {
	value, err := uuid.Parse(p.pathFunc(r, name))
	if err != nil {
		p.logger.Warn(r.Context(), "PathParamUUID", "error", err)

		return uuid.Nil
	}

	return value
}

// FilterUUID - возвращает UUID поступивший из внешнего запроса.
// Если ключ key не найден или возникнет ошибка, то возвращается uuid.Nil.
func (p *UUID) FilterUUID(r *http.Request, key string) uuid.UUID {
	value, err := mrreq.ParseUUID(r.URL.Query(), key, false)
	if err != nil {
		p.logger.Warn(r.Context(), "FilterUUID", "error", err)

		return uuid.Nil
	}

	return value
}
