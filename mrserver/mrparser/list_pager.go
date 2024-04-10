package mrparser

import (
	"net/http"

	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/mrreq"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	ListPager struct {
		paramNamePageIndex string
		paramNamePageSize  string
		pageSizeMax        uint64
		pageSizeDefault    uint64
	}

	ListPagerOptions struct {
		ParamNamePageIndex string
		ParamNamePageSize  string
		PageSizeMax        uint64
		PageSizeDefault    uint64
	}
)

var (
	// Make sure the ListPager conforms with the mrserver.RequestParserListPager interface
	_ mrserver.RequestParserListPager = (*ListPager)(nil)
)

func NewListPager(opts ListPagerOptions) *ListPager {
	lp := ListPager{
		paramNamePageIndex: "pageIndex",
		paramNamePageSize:  "pageSize",
		pageSizeMax:        1000,
		pageSizeDefault:    10,
	}

	if opts.ParamNamePageIndex != "" {
		lp.paramNamePageIndex = opts.ParamNamePageIndex
	}

	if opts.ParamNamePageSize != "" {
		lp.paramNamePageSize = opts.ParamNamePageSize
	}

	if opts.PageSizeMax > 0 {
		lp.pageSizeMax = opts.PageSizeMax
	}

	if opts.PageSizeDefault > 0 {
		lp.pageSizeDefault = opts.PageSizeDefault
	}

	return &lp
}

func (p *ListPager) PageParams(r *http.Request) mrtype.PageParams {
	value, err := mrreq.ParsePageParams(
		r,
		p.paramNamePageIndex,
		p.paramNamePageSize,
	)

	if err != nil || value.Size < 1 || value.Size > p.pageSizeMax {
		if err != nil {
			mrlog.Ctx(r.Context()).Warn().Err(err).Send()
		}

		return mrtype.PageParams{
			Size: p.pageSizeDefault,
		}
	}

	return value
}
