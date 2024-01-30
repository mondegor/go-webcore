package mrparser

import (
	"net/http"
	"time"

	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/mrreq"
)

type (
	DateTime struct {
	}
)

// Make sure the DateTime conforms with the mrserver.RequestParserDateTime interface
var _ mrserver.RequestParserDateTime = (*DateTime)(nil)

func NewDateTime() *DateTime {
	return &DateTime{}
}

func (p *DateTime) FilterDateTime(r *http.Request, key string) time.Time {
	value, err := mrreq.ParseDateTime(r, key, false)

	if err != nil {
		mrlog.Ctx(r.Context()).Warn().Err(err)
		return time.Time{}
	}

	return value
}
