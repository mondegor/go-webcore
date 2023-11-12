package mrview

import (
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	ListSorter interface {
		CheckField(name string) bool
		DefaultSort() mrtype.SortParams
	}
)
