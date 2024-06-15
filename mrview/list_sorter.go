package mrview

import (
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	// ListSorter - интерфейс для проверки полей, которые могут участвовать в сортировке.
	ListSorter interface {
		CheckField(name string) bool
		DefaultSort() mrtype.SortParams
	}
)
