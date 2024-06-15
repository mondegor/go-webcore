package mrtype

import "github.com/mondegor/go-webcore/mrenum"

type (
	// KeyInt32 - comment type.
	KeyInt32 int32

	// KeyInt64 - comment type.
	KeyInt64 int64

	// KeyString - comment type.
	KeyString string

	// RangeInt64 - comment struct.
	RangeInt64 struct {
		Min int64
		Max int64
	}

	// PageParams - comment struct.
	PageParams struct {
		Index uint64 // pageIndex
		Size  uint64 // pageSize
	}

	// SortParams - comment struct.
	SortParams struct {
		FieldName string               // sortField
		Direction mrenum.SortDirection // sortDirection
	}
)
