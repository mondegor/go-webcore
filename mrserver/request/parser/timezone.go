package parser

import (
	"context"
	"net/http"
	"time"

	"github.com/mondegor/go-core/mrlog"
	"github.com/mondegor/go-core/util/timezone"

	"github.com/mondegor/go-webcore/mrserver"
)

type (
	// TimeZone - определяет часовой пояс запроса.
	TimeZone struct {
		locations   locationList
		logger      mrlog.Logger
		paramNameTZ string
	}

	// locationList - предоставляет доступ к предзагруженным часовым поясам приложения:
	// к поясу по имени (строгая проверка принадлежности списку), к подбору пояса
	// по смещению и к поясу по умолчанию.
	locationList interface {
		LocationByName(value string) (*time.Location, error)
		NameByOffset(offset time.Duration, isDST bool) (name string, ok bool)
		Default() *time.Location
	}
)

// NewTimeZone - создаёт объект TimeZone.
// Параметры:
//   - locations - список часовых поясов, зарегистрированных приложением при старте;
//   - logger - логгер;
//   - paramNameTZ - имя query-параметра разового переопределения пояса (например: "tz").
func NewTimeZone(
	locations locationList,
	logger mrlog.Logger,
	paramNameTZ string,
) *TimeZone {
	return &TimeZone{
		locations:   locations,
		logger:      logger,
		paramNameTZ: paramNameTZ,
	}
}

// Location - определяет часовой пояс запроса; порядок источников и контракт описаны
// в request.ParserTimeZone. Если пояс не удалось определить ни по одному источнику,
// возвращается пояс по умолчанию - первый пояс списка, зарегистрированного приложением при старте.
//
// Первые два источника (query-параметр и X-Internal-Time-Zone) строгие и потому разбираются
// одним проходом: негодное значение молча отбрасывается, и разбор продолжается со следующего
// источника. Логирования здесь нет намеренно: годные значения клиент получает от самого
// приложения, а промах на внутреннем заголовке означает пояс, который приложение больше
// не предоставляет, - и то, и другое штатно откатывается ниже.
//
// X-Accept-Time-Zone разбирается терпимее, так как приходит от браузера, который о списке поясов
// приложения ничего не знает: незарегистрированное имя не отбрасывает заголовок целиком,
// а передаёт разбор параметрам offset и dst, по которым подбирается ближайший
// из поясов приложения (см. LocationList.NameByOffset).
//
// Набор доступных поясов ограничен списком, который приложение зарегистрировало при старте,
// поэтому обращения к базе часовых поясов в процессе обработки запроса не происходит.
func (p *TimeZone) Location(r *http.Request) *time.Location {
	ctx := r.Context()

	// строгие источники в порядке приоритета
	for _, name := range [...]string{
		r.URL.Query().Get(p.paramNameTZ),
		r.Header.Get(mrserver.HeaderKeyInternalTimeZone),
	} {
		if name == "" {
			continue
		}

		if loc, ok := p.locationByName(name); ok {
			return loc
		}
	}

	if value := r.Header.Get(mrserver.HeaderKeyAcceptTimeZone); value != "" {
		if loc, ok := p.locationByHeader(ctx, value); ok {
			return loc
		}
	}

	return p.locations.Default()
}

// locationByName - возвращает пояс по имени, если оно есть в списке приложения.
//
// Проверка строгая и целиком на стороне LocationByName: он принимает имя только при точном
// совпадении с поясом списка, а "Local" (пояс процесса) наружу не отдаёт, поэтому промах
// здесь - штатный случай (незарегистрированное имя), а не сбой, и не логируется.
func (p *TimeZone) locationByName(name string) (*time.Location, bool) {
	loc, err := p.locations.LocationByName(name)
	if err != nil {
		return nil, false
	}

	return loc, true
}

// locationByHeader - определяет пояс по значению заголовка: сначала по имени,
// затем подбором по смещению. Сообщает false, если не удалось ни то, ни другое.
func (p *TimeZone) locationByHeader(ctx context.Context, value string) (*time.Location, bool) {
	header, err := timezone.ParseAcceptTimeZone(value)
	if err != nil {
		p.logger.Warn(
			ctx,
			"Header is incorrect",
			"header", mrserver.HeaderKeyAcceptTimeZone,
			"timezone", value,
		)

		return nil, false
	}

	if header.Name != "" {
		if loc, found := p.locationByName(header.Name); found {
			return loc, true
		}
	}

	// имя не подошло, поэтому в дело вступает запасной путь: по смещению
	// подбирается ближайший из поясов приложения. Промах здесь не логируется -
	// клиент в поясе, которого приложение не предоставляет, штатный случай
	if header.HasOffset {
		if name, found := p.locations.NameByOffset(header.Offset, header.IsDST); found {
			return p.locationByName(name)
		}
	}

	return nil, false
}
