package request

import (
	"context"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/mondegor/go-sysmess/mrstatus/itemstatus"
	"github.com/mondegor/go-sysmess/mrtype"

	"github.com/mondegor/go-webcore/mrcore"
)

type (
	// ParserDecode - преобразователь данных в указанную go структуру.
	ParserDecode interface {
		ParseToStruct(ctx context.Context, content io.Reader, structPointer any) error
	}

	// ParserString - парсер данных запроса для преобразования их в string.
	ParserString interface {
		PathParamString(r *http.Request, name string) string
		RawParamString(r *http.Request, key string) *string // returns nil if the param not found
		FilterString(r *http.Request, key string) string
	}

	// ParserValidate - парсер данных запроса для преобразования их в go структуру.
	ParserValidate interface {
		Validate(r *http.Request, structPointer any) error
		ValidateContent(ctx context.Context, content []byte, structPointer any) error
	}

	// ParserInt64 - парсер данных запроса для преобразования их в int64.
	ParserInt64 interface {
		FilterInt64(r *http.Request, key string) int64
		FilterRangeInt64(r *http.Request, key string) mrtype.RangeInt64
		FilterInt64List(r *http.Request, key string) []int64
	}

	// ParserUint64 - парсер данных запроса для преобразования их в mrtype.Uint64.
	ParserUint64 interface {
		PathParamUint64(r *http.Request, name string) uint64
		FilterUint64(r *http.Request, key string) uint64
		FilterUint64List(r *http.Request, key string) []uint64
	}

	// ParserFloat64 - парсер данных запроса для преобразования их во float64.
	ParserFloat64 interface {
		FilterFloat64(r *http.Request, key string) float64
		FilterRangeFloat64(r *http.Request, key string) mrtype.RangeFloat64
	}

	// ParserBool - парсер данных запроса для преобразования их в *bool.
	ParserBool interface {
		FilterNullableBool(r *http.Request, key string) *bool
	}

	// ParserDateTime - парсер данных запроса для преобразования их в time.Time.
	ParserDateTime interface {
		FilterDateTime(r *http.Request, key string) time.Time
	}

	// ParserUUID - парсер данных запроса для преобразования их в uuid.UUID.
	ParserUUID interface {
		PathParamUUID(r *http.Request, name string) uuid.UUID
		FilterUUID(r *http.Request, key string) uuid.UUID
	}

	// ParserFile - парсер данных запроса для преобразования их в файловую структуру.
	ParserFile interface {
		FormFile(r *http.Request, key string) (mrtype.File, error)
		FormFileContent(r *http.Request, key string) (mrtype.FileContent, error)
		FormFiles(r *http.Request, key string) ([]mrtype.FileHeader, error)
	}

	// ParserImage - парсер данных запроса для преобразования их в файловую структуру изображения.
	ParserImage interface {
		FormImage(r *http.Request, key string) (mrtype.Image, error)
		FormImageContent(r *http.Request, key string) (mrtype.ImageContent, error)
		FormImages(r *http.Request, key string) ([]mrtype.ImageHeader, error)
	}

	// ParserListSorter - парсер данных запроса для преобразования их в mrtype.SortParams.
	ParserListSorter interface {
		SortParams(r *http.Request, sorter mrtype.ListSorter) mrtype.SortParams
	}

	// ParserListPager - парсер данных запроса для преобразования их в mrtype.PageParams.
	ParserListPager interface {
		PageParams(r *http.Request) mrtype.PageParams
	}

	// ParserItemStatus - парсер данных запроса для преобразования их в []mrenum.ItemStatus.
	ParserItemStatus interface {
		FilterStatusList(r *http.Request, key string) []itemstatus.Enum
	}

	// ParserEnumList - парсер данных запроса для преобразования их в []mrenum.ItemStatus.
	ParserEnumList[T ~uint8] interface {
		FilterEnumList(r *http.Request, key string) []T
	}

	// ParserClientIP - comment interface.
	ParserClientIP interface {
		RealIP(r *http.Request) net.IP
		DetailedIP(r *http.Request) mrtype.DetailedIP
	}

	// ParserLocale - comment interface.
	ParserLocale interface {
		Language(r *http.Request) string
		Localizer(r *http.Request) mrcore.Localizer
	}

	// ParserUser - comment interface.
	ParserUser interface {
		UserID(r *http.Request) uuid.UUID
		UserAndGroup(r *http.Request) (userID uuid.UUID, group string)
	}
)
