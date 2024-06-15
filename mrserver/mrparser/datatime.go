package mrparser

import (
	"net/http"
	"time"

	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/mrreq"
)

type (
	// DateTime - comment struct.
	DateTime struct{}
)

// Make sure the DateTime conforms with the mrserver.RequestParserDateTime interface.
var _ mrserver.RequestParserDateTime = (*DateTime)(nil)

// NewDateTime - создаёт объект DateTime.
func NewDateTime() *DateTime {
	return &DateTime{}
}

// FilterDateTime - comment method.
func (p *DateTime) FilterDateTime(r *http.Request, key string) time.Time {
	value, err := mrreq.ParseDateTime(r, key, false)
	if err != nil {
		mrlog.Ctx(r.Context()).Warn().Err(err).Send()

		return time.Time{}
	}

	return value
}
