package mrserver

import (
	"context"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/mondegor/go-sysmess/mrtype"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrenum"
)

type (
	// RequestDecoder - преобразователь данных в указанную go структуру.
	RequestDecoder interface {
		ParseToStruct(ctx context.Context, content io.Reader, structPointer any) error
	}

	// RequestParserString - парсер данных запроса для преобразования их в string.
	RequestParserString interface {
		PathParamString(r *http.Request, name string) string
		RawParamString(r *http.Request, key string) *string // returns nil if the param not found
		FilterString(r *http.Request, key string) string
	}

	// RequestParserValidate - парсер данных запроса для преобразования их в go структуру.
	RequestParserValidate interface {
		Validate(r *http.Request, structPointer any) error
		ValidateContent(ctx context.Context, content []byte, structPointer any) error
	}

	// RequestParserInt64 - парсер данных запроса для преобразования их в int64.
	RequestParserInt64 interface {
		FilterInt64(r *http.Request, key string) int64
		FilterRangeInt64(r *http.Request, key string) mrtype.RangeInt64
		FilterInt64List(r *http.Request, key string) []int64
	}

	// RequestParserUint64 - парсер данных запроса для преобразования их в mrtype.Uint64.
	RequestParserUint64 interface {
		PathParamUint64(r *http.Request, name string) uint64
		FilterUint64(r *http.Request, key string) uint64
		FilterUint64List(r *http.Request, key string) []uint64
	}

	// RequestParserFloat64 - парсер данных запроса для преобразования их во float64.
	RequestParserFloat64 interface {
		FilterFloat64(r *http.Request, key string) float64
		FilterRangeFloat64(r *http.Request, key string) mrtype.RangeFloat64
	}

	// RequestParserBool - парсер данных запроса для преобразования их в *bool.
	RequestParserBool interface {
		FilterNullableBool(r *http.Request, key string) *bool
	}

	// RequestParserDateTime - парсер данных запроса для преобразования их в time.Time.
	RequestParserDateTime interface {
		FilterDateTime(r *http.Request, key string) time.Time
	}

	// RequestParserUUID - парсер данных запроса для преобразования их в uuid.UUID.
	RequestParserUUID interface {
		PathParamUUID(r *http.Request, name string) uuid.UUID
		FilterUUID(r *http.Request, key string) uuid.UUID
	}

	// RequestParserFile - парсер данных запроса для преобразования их в файловую структуру.
	RequestParserFile interface {
		FormFile(r *http.Request, key string) (mrtype.File, error)
		FormFileContent(r *http.Request, key string) (mrtype.FileContent, error)
		FormFiles(r *http.Request, key string) ([]mrtype.FileHeader, error)
	}

	// RequestParserImage - парсер данных запроса для преобразования их в файловую структуру изображения.
	RequestParserImage interface {
		FormImage(r *http.Request, key string) (mrtype.Image, error)
		FormImageContent(r *http.Request, key string) (mrtype.ImageContent, error)
		FormImages(r *http.Request, key string) ([]mrtype.ImageHeader, error)
	}

	// RequestParserListSorter - парсер данных запроса для преобразования их в mrtype.SortParams.
	RequestParserListSorter interface {
		SortParams(r *http.Request, sorter mrtype.ListSorter) mrtype.SortParams
	}

	// RequestParserListPager - парсер данных запроса для преобразования их в mrtype.PageParams.
	RequestParserListPager interface {
		PageParams(r *http.Request) mrtype.PageParams
	}

	// RequestParserItemStatus - парсер данных запроса для преобразования их в []mrenum.ItemStatus.
	RequestParserItemStatus interface {
		FilterStatusList(r *http.Request, key string) []mrenum.ItemStatus
	}

	// RequestParserClientIP - comment interface.
	RequestParserClientIP interface {
		RealIP(r *http.Request) net.IP
		DetailedIP(r *http.Request) mrtype.DetailedIP
	}

	// RequestParserLocale - comment interface.
	RequestParserLocale interface {
		Language(r *http.Request) string
		Localizer(r *http.Request) mrcore.Localizer
	}

	// RequestParserUser - comment interface.
	RequestParserUser interface {
		UserID(r *http.Request) uuid.UUID
		UserAndGroup(r *http.Request) (userID uuid.UUID, group string)
	}

	// RequestParserParamFunc - функция для парсинга URL для извлечения из него параметров.
	RequestParserParamFunc func(r *http.Request, key string) string
)
