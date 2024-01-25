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

	RequestParserString interface {
		PathParamString(r *http.Request, name string) string
		// RawParamString - returns nil if the param not found
		RawParamString(r *http.Request, key string) *string
		FilterString(r *http.Request, key string) string
	}

	RequestParserValidate interface {
		Validate(r *http.Request, structPointer any) error
	}

	RequestParserInt64 interface {
		PathParamInt64(r *http.Request, name string) int64
		FilterInt64(r *http.Request, key string) int64
		FilterRangeInt64(r *http.Request, key string) mrtype.RangeInt64
		FilterInt64List(r *http.Request, key string) []int64
	}

	RequestParserKeyInt32 interface {
		PathKeyInt32(r *http.Request, name string) mrtype.KeyInt32
		FilterKeyInt32(r *http.Request, key string) mrtype.KeyInt32
		FilterKeyInt32List(r *http.Request, key string) []mrtype.KeyInt32
	}

	RequestParserBool interface {
		FilterNullableBool(r *http.Request, key string) *bool
	}

	RequestParserDateTime interface {
		FilterDateTime(r *http.Request, key string) time.Time
	}

	RequestParserUUID interface {
		PathParamUUID(r *http.Request, name string) uuid.UUID
		FilterUUID(r *http.Request, key string) uuid.UUID
	}

	RequestParserFile interface {
		FormFile(r *http.Request, key string) (mrtype.File, error)
		FormFileContent(r *http.Request, key string) (mrtype.FileContent, error)
	}

	RequestParserImage interface {
		FormImage(r *http.Request, key string) (mrtype.Image, error)
		FormImageContent(r *http.Request, key string) (mrtype.ImageContent, error)
	}

	RequestParserSortPage interface {
		SortParams(r *http.Request, sorter mrview.ListSorter) mrtype.SortParams
		PageParams(r *http.Request) mrtype.PageParams
	}

	RequestParserItemStatus interface {
		FilterStatusList(r *http.Request, key string) []mrenum.ItemStatus
	}

	RequestParserParamFunc func(r *http.Request, key string) string
)
