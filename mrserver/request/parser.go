package request

import (
	"context"
	"io"
	"net/http"
	"net/netip"
	"time"

	"github.com/google/uuid"
	"github.com/mondegor/go-core/mrmodel/media"
	"github.com/mondegor/go-core/mrtype"
	"github.com/mondegor/go-core/mrworkflow/itemstatus"

	"github.com/mondegor/go-webcore/mrcore"
)

type (
	// ParserDecode - преобразует данные запроса в Go-структуру.
	// Используется для парсинга тела HTTP-запроса (например: JSON) в структуру.
	ParserDecode interface {
		ParseToStruct(ctx context.Context, content io.Reader, structPointer any) error
	}

	// ParserString - извлекает строковые данные из HTTP-запроса.
	//
	// Поддерживает извлечение из:
	//  - PathParamString - параметры пути (например: /users/{id});
	//  - RawParamString - query-параметры (например: ?name=value);
	//  - FilterString/FilterStringList - фильтры из query-параметров.
	ParserString interface {
		PathParamString(r *http.Request, name string) string
		RawParamString(r *http.Request, key string) *string // returns nil if the param not found
		FilterString(r *http.Request, key string) string
		FilterStringList(r *http.Request, key string) []string
	}

	// ParserValidate - валидирует данные запроса, преобразованные в Go-структуру.
	// Использует теги валидации в структуре (например: `validate:"required"`).
	ParserValidate interface {
		Validate(r *http.Request, structPointer any) error
		ValidateContent(ctx context.Context, content []byte, structPointer any) error
	}

	// ParserInt64 - извлекает целочисленные значения int64 из query-параметров.
	// Поддерживает: одиночные значения, диапазоны (-min/-max), списки.
	ParserInt64 interface {
		FilterInt64(r *http.Request, key string) int64
		FilterRangeInt64(r *http.Request, key string) mrtype.RangeInt64
		FilterInt64List(r *http.Request, key string) []int64
	}

	// ParserUint64 - извлекает целочисленные значения uint64 из запроса.
	// Поддерживает извлечение из URL-пути и query-параметров.
	ParserUint64 interface {
		PathParamUint64(r *http.Request, name string) uint64
		FilterUint64(r *http.Request, key string) uint64
		FilterUint64List(r *http.Request, key string) []uint64
	}

	// ParserFloat64 - извлекает значения с плавающей точкой float64 из query-параметров.
	ParserFloat64 interface {
		FilterFloat64(r *http.Request, key string) float64
		FilterRangeFloat64(r *http.Request, key string) mrtype.RangeFloat64
	}

	// ParserBool - извлекает логические значения из query-параметров.
	ParserBool interface {
		FilterNullableBool(r *http.Request, key string) *bool
	}

	// ParserDateTime - извлекает значения даты и времени time.Time из query-параметров.
	ParserDateTime interface {
		FilterDateTime(r *http.Request, key string) time.Time
	}

	// ParserUUID - извлекает UUID-значения из запроса.
	ParserUUID interface {
		PathParamUUID(r *http.Request, name string) uuid.UUID
		FilterUUID(r *http.Request, key string) uuid.UUID
	}

	// ParserFile - извлекает файлы и их метаданные из multipart-формы запроса.
	ParserFile interface {
		FormFile(r *http.Request, key string) (media.File, error)
		FormFileContent(r *http.Request, key string) (media.FileContent, error)
		FormFiles(r *http.Request, key string) ([]media.FileHeader, error)
	}

	// ParserImage - извлекает файлы изображений и их метаданные из multipart-формы запроса.
	ParserImage interface {
		FormImage(r *http.Request, key string) (media.Image, error)
		FormImageContent(r *http.Request, key string) (media.ImageContent, error)
		FormImages(r *http.Request, key string) ([]media.ImageHeader, error)
	}

	// ParserListSorter - извлекает параметры сортировки из query-параметров запроса.
	ParserListSorter interface {
		SortParams(r *http.Request, sorter mrtype.ListSorter) mrtype.SortParams
	}

	// ParserListPager - извлекает параметры постраничной пагинации из запроса.
	ParserListPager interface {
		PageParams(r *http.Request) mrtype.PageParams
	}

	// ParserListCursor - извлекает параметры курсорной пагинации из запроса.
	ParserListCursor interface {
		CursorParams(r *http.Request) mrtype.CursorParams
	}

	// ParserItemStatus - извлекает список статусов элементов из query-параметра.
	ParserItemStatus interface {
		FilterStatusList(r *http.Request, key string) []itemstatus.Enum
	}

	// ParserEnumList - извлекает список перечислений заданного типа из query-параметра.
	ParserEnumList[T ~uint8] interface {
		FilterEnumList(r *http.Request, key string) []T
	}

	// ParserClientIP - парсер для получения IP-адреса клиента из запроса.
	ParserClientIP interface {
		RealIP(r *http.Request) netip.Addr
		DetailedIP(r *http.Request) mrtype.DetailedIP
	}

	// ParserLocale - парсер для определения локали и языка из запроса.
	// Определяет язык из query-параметра или заголовка Accept-Language.
	ParserLocale interface {
		Language(r *http.Request) string
		Localizer(r *http.Request) mrcore.Localizer
	}

	// ParserUser - парсер для получения информации о пользователе из запроса.
	// Извлекает данные из внутренних заголовков X-Internal-UserId-Group,
	// X-Internal-Session-Id и X-Internal-Time-Zone, установленных middleware.
	ParserUser interface {
		UserID(r *http.Request) uuid.UUID
		UserAndGroup(r *http.Request) (userID uuid.UUID, group string)
		Location(r *http.Request) *time.Location
		SessionID(r *http.Request) string
	}
)
