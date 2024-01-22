package mrparser

import (
	"net/http"

	"github.com/mondegor/go-webcore/mrctx"
	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/mrreq"
	"github.com/mondegor/go-webcore/mrtype"
	"github.com/mondegor/go-webcore/mrview"
)

type (
	SortPage struct {
		paramNameSortField     string
		paramNameSortDirection string
		paramNamePageIndex     string
		paramNamePageSize      string

		pageSizeMax     uint64
		pageSizeDefault uint64
	}
)

// Make sure the SortPage conforms with the mrserver.RequestParserSortPage interface
var _ mrserver.RequestParserSortPage = (*SortPage)(nil)

func NewSortPage() *SortPage {
	return &SortPage{
		paramNamePageIndex:     "pageIndex",
		paramNamePageSize:      "pageSize",
		paramNameSortField:     "sortField",
		paramNameSortDirection: "sortDirection",
		pageSizeMax:            1000,
		pageSizeDefault:        10,
	}
}

func (p *SortPage) SortParams(r *http.Request, sorter mrview.ListSorter) mrtype.SortParams {
	value, err := mrreq.ParseSortParams(
		r,
		p.paramNameSortField,
		p.paramNameSortDirection,
	)

	if err != nil {
		mrctx.Logger(r.Context()).Warn(err)
	}

	if value.FieldName == "" {
		return sorter.DefaultSort()
	}

	if !sorter.CheckField(value.FieldName) {
		mrctx.Logger(r.Context()).Warning("sort field '%s' is not registered", value.FieldName)
		return sorter.DefaultSort()
	}

	return value
}

func (p *SortPage) PageParams(r *http.Request) mrtype.PageParams {
	value, err := mrreq.ParsePageParams(
		r,
		p.paramNamePageIndex,
		p.paramNamePageSize,
	)

	// :TODO: вынести параметры p.pageSizeMax и p.pageSizeDefault по аналогии с SortParams
	if err != nil || value.Size < 1 || value.Size > p.pageSizeMax {
		if err != nil {
			mrctx.Logger(r.Context()).Warn(err)
		}

		return mrtype.PageParams{
			Size: p.pageSizeDefault,
		}
	}

	return value
}
