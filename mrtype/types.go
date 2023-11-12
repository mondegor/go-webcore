package mrtype

import "github.com/mondegor/go-webcore/mrenum"

type (
    KeyInt32 int32
    KeyInt64 int64
    KeyString string

    RangeInt64 struct {
        Min int64
        Max int64
    }

    PageParams struct {
        Index uint64 // pageIndex
        Size  uint64 // pageSize
    }

    SortParams struct {
        FieldName string // sortField
        Direction mrenum.SortDirection // sortDirection
    }
)
