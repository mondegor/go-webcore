package parser

import (
	"net/http"

	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-sysmess/mrtype"
	"github.com/mondegor/go-sysmess/mrtype/parse"
)

type (
	// ListSorter - парсер параметров для сортировки списка элементов.
	ListSorter struct {
		logger                 mrlog.Logger
		paramNameSortColumn    string
		paramNameSortDirection string
	}

	// ListSorterOptions - опции для создания ListSorter.
	ListSorterOptions struct {
		ParamNameSortColumn    string
		ParamNameSortDirection string
	}
)

// NewListSorter - создаёт объект ListSorter.
func NewListSorter(logger mrlog.Logger, opts ListSorterOptions) *ListSorter {
	ls := ListSorter{
		logger:                 logger,
		paramNameSortColumn:    "sortColumn",
		paramNameSortDirection: "sortDirection",
	}

	if opts.ParamNameSortColumn != "" {
		ls.paramNameSortColumn = opts.ParamNameSortColumn
	}

	if opts.ParamNameSortDirection != "" {
		ls.paramNameSortDirection = opts.ParamNameSortDirection
	}

	return &ls
}

// SortParams - возвращает распарсенные параметры сортировки списка элементов.
func (p *ListSorter) SortParams(r *http.Request, sorter mrtype.ListSorter) mrtype.SortParams {
	value, err := parse.SortParams(
		r.URL.Query().Get(p.paramNameSortColumn),
		r.URL.Query().Get(p.paramNameSortDirection),
	)
	if err != nil {
		p.logger.Warn(
			r.Context(), "SortParams",
			"column", p.paramNameSortColumn,
			"direction", p.paramNameSortDirection,
			"error", err,
		)
	}

	if value.Column == "" {
		return sorter.DefaultSort()
	}

	if !sorter.HasColumn(value.Column) {
		p.logger.Warn(r.Context(), "sort column is not registered", "column", value.Column)

		return sorter.DefaultSort()
	}

	return value
}
