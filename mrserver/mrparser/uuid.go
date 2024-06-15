package mrparser

import (
	"net/http"

	"github.com/google/uuid"

	"github.com/mondegor/go-webcore/mrlog"

	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/mrreq"
)

type (
	// UUID - comment struct.
	UUID struct {
		pathFunc mrserver.RequestParserParamFunc
	}
)

// Make sure the UUID conforms with the mrserver.RequestParserUUID interface.
var _ mrserver.RequestParserUUID = (*UUID)(nil)

// NewUUID - создаёт объект UUID.
func NewUUID(pathFunc mrserver.RequestParserParamFunc) *UUID {
	return &UUID{
		pathFunc: pathFunc,
	}
}

// PathParamUUID - comment method.
func (p *UUID) PathParamUUID(r *http.Request, name string) uuid.UUID {
	value, err := uuid.Parse(p.pathFunc(r, name))
	if err != nil {
		mrlog.Ctx(r.Context()).Warn().Err(err).Send()

		return uuid.Nil
	}

	return value
}

// FilterUUID - comment method.
func (p *UUID) FilterUUID(r *http.Request, key string) uuid.UUID {
	value, err := mrreq.ParseUUID(r, key, false)
	if err != nil {
		mrlog.Ctx(r.Context()).Warn().Err(err).Send()

		return uuid.Nil
	}

	return value
}
