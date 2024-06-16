package mrtype

import "github.com/mondegor/go-webcore/mrenum"

type (
	// KeyInt32 - целочисленный ID записи.
	KeyInt32 int32

	// KeyInt64 - целочисленный ID записи.
	KeyInt64 int64

	// KeyString - строковый ID записи.
	KeyString string

	// RangeInt64 - comment struct.
	RangeInt64 struct {
		Min int64
		Max int64
	}

	// PageParams - параметры для выборки части списка элементов.
	PageParams struct {
		Index uint64 // pageIndex
		Size  uint64 // pageSize
	}

	// SortParams - параметры для сортировки списка элементов по указанному полю.
	SortParams struct {
		FieldName string               // sortField
		Direction mrenum.SortDirection // sortDirection
	}
)
