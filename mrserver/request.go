package mrserver

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"
	"github.com/mondegor/go-webcore/mrview"
)

type (
	RequestDecoder interface {
		ParseToStruct(r *http.Request, structPointer any) error
	}

	RequestParserPath interface {
		PathParam(r *http.Request, name string) string
	}

	RequestParserPathFunc func(r *http.Request, name string) string

	RequestParser interface {
		PathParamString(r *http.Request, name string) string
		PathParamInt64(r *http.Request, name string) int64

		// RawQueryParam - returns nil if the param not found
		RawQueryParam(r *http.Request, key string) *string

		FilterString(r *http.Request, key string) string
		FilterNullableBool(r *http.Request, key string) *bool
		FilterInt64(r *http.Request, key string) int64
		FilterRangeInt64(r *http.Request, key string) mrtype.RangeInt64
		FilterInt64List(r *http.Request, key string) []int64
		FilterDateTime(r *http.Request, key string) time.Time
	}

	RequestParserItemStatus interface {
		FilterStatusList(r *http.Request, key string) []mrenum.ItemStatus
	}

	RequestParserKeyInt32 interface {
		PathKeyInt32(r *http.Request, name string) mrtype.KeyInt32
		FilterKeyInt32(r *http.Request, key string) mrtype.KeyInt32
		FilterKeyInt32List(r *http.Request, key string) []mrtype.KeyInt32
	}

	RequestParserSortPage interface {
		SortParams(r *http.Request, sorter mrview.ListSorter) mrtype.SortParams
		PageParams(r *http.Request) mrtype.PageParams
	}

	RequestParserUUID interface {
		PathParamUUID(r *http.Request, name string) uuid.UUID
		FilterUUID(r *http.Request, key string) uuid.UUID
	}

	RequestParserValidate interface {
		Validate(r *http.Request, structPointer any) error
	}
)

func (f RequestParserPathFunc) PathParam(r *http.Request, name string) string {
	return f(r, name)
}
