package mrparser

import (
	"net/http"

	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/mrreq"
	"github.com/mondegor/go-webcore/mrtype"
	"github.com/mondegor/go-webcore/mrview"
)

type (
	ListSorter struct {
		paramNameSortField     string
		paramNameSortDirection string
	}

	ListSorterOptions struct {
		ParamNameSortField     string
		ParamNameSortDirection string
	}
)

// Make sure the ListSorter conforms with the mrserver.RequestParserListSorter interface
var _ mrserver.RequestParserListSorter = (*ListSorter)(nil)

func NewListSorter(opts ListSorterOptions) *ListSorter {
	ls := ListSorter{
		paramNameSortField:     "sortField",
		paramNameSortDirection: "sortDirection",
	}

	if opts.ParamNameSortField != "" {
		ls.paramNameSortField = opts.ParamNameSortField
	}

	if opts.ParamNameSortDirection != "" {
		ls.paramNameSortDirection = opts.ParamNameSortDirection
	}

	return &ls
}

func (p *ListSorter) SortParams(r *http.Request, sorter mrview.ListSorter) mrtype.SortParams {
	value, err := mrreq.ParseSortParams(
		r,
		p.paramNameSortField,
		p.paramNameSortDirection,
	)

	if err != nil {
		mrlog.Ctx(r.Context()).Warn().Err(err).Send()
	}

	if value.FieldName == "" {
		return sorter.DefaultSort()
	}

	if !sorter.CheckField(value.FieldName) {
		mrlog.Ctx(r.Context()).Warn().Msgf("sort field '%s' is not registered", value.FieldName)
		return sorter.DefaultSort()
	}

	return value
}
