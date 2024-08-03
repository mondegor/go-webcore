package mrparser

import (
	"net/http"

	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/mrreq"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	// Float64 - comment struct.
	Float64 struct{}
)

// Make sure the Float64 conforms with the mrserver.RequestParserFloat64 interface.
var _ mrserver.RequestParserFloat64 = (*Float64)(nil)

// NewFloat64 - создаёт объект Float64.
func NewFloat64() *Float64 {
	return &Float64{}
}

// FilterFloat64 - comment method.
func (p *Float64) FilterFloat64(r *http.Request, key string) float64 {
	value, err := mrreq.ParseFloat64(r, key, false)
	if err != nil {
		mrlog.Ctx(r.Context()).Warn().Err(err).Send()

		return 0
	}

	return value
}

// FilterRangeFloat64 - comment method.
func (p *Float64) FilterRangeFloat64(r *http.Request, key string) mrtype.RangeFloat64 {
	value, err := mrreq.ParseRangeFloat64(r, key)
	if err != nil {
		mrlog.Ctx(r.Context()).Warn().Err(err).Send()

		return mrtype.RangeFloat64{}
	}

	return value
}
