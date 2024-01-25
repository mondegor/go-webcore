package mrparser

import (
	"net/http"

	"github.com/mondegor/go-webcore/mrctx"
	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/mrreq"
)

type (
	Bool struct {
	}
)

// Make sure the Bool conforms with the mrserver.RequestParserBool interface
var _ mrserver.RequestParserBool = (*Bool)(nil)

func NewBool() *Bool {
	return &Bool{}
}

func (p *Bool) FilterNullableBool(r *http.Request, key string) *bool {
	value, err := mrreq.ParseNullableBool(r, key)

	if err != nil {
		mrctx.Logger(r.Context()).Warn(err)
		return nil
	}

	return value
}
