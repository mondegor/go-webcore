package mrparser

import (
	"net/http"

	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrserver/mrreq"
	"github.com/mondegor/go-webcore/mrtype"
	"github.com/mondegor/go-webcore/mrview"
)

type (
	// ListSorter - парсер параметров для сортировки списка элементов.
	ListSorter struct {
		paramNameSortField     string
		paramNameSortDirection string
	}

	// ListSorterOptions - опции для создания ListSorter.
	ListSorterOptions struct {
		ParamNameSortField     string
		ParamNameSortDirection string
	}
)

// NewListSorter - создаёт объект ListSorter.
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

// SortParams - возвращает распарсенные параметры сортировки списка элементов.
func (p *ListSorter) SortParams(r *http.Request, sorter mrview.ListSorter) mrtype.SortParams {
	value, err := mrreq.ParseSortParams(
		r.URL.Query(),
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
