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
func NewListSorter(logger mrlog.Logger, opts ListSorterOptions) *ListSorter {
	ls := ListSorter{
		logger:                 logger,
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
func (p *ListSorter) SortParams(r *http.Request, sorter mrtype.ListSorter) mrtype.SortParams {
	value, err := parse.SortParams(
		r.URL.Query().Get(p.paramNameSortField),
		r.URL.Query().Get(p.paramNameSortDirection),
	)
	if err != nil {
		p.logger.Warn(
			r.Context(), "SortParams",
			"field", p.paramNameSortField,
			"direction", p.paramNameSortDirection,
			"error", err,
		)
	}

	if value.FieldName == "" {
		return sorter.DefaultSort()
	}

	if !sorter.HasField(value.FieldName) {
		p.logger.Warn(r.Context(), "sort field is not registered", "field", value.FieldName)

		return sorter.DefaultSort()
	}

	return value
}
