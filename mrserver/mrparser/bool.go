package mrparser

import (
	"net/http"

	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/mrreq"
)

type (
	// Bool - comment struct.
	Bool struct{}
)

// Make sure the Bool conforms with the mrserver.RequestParserBool interface.
var _ mrserver.RequestParserBool = (*Bool)(nil)

// NewBool - создаёт объект Bool.
func NewBool() *Bool {
	return &Bool{}
}

// FilterNullableBool - comment method.
func (p *Bool) FilterNullableBool(r *http.Request, key string) *bool {
	value, err := mrreq.ParseNullableBool(r, key)
	if err != nil {
		mrlog.Ctx(r.Context()).Warn().Err(err).Send()

		return nil
	}

	return value
}
