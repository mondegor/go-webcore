package parser_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/mondegor/go-core/mrlog"
	"github.com/mondegor/go-core/util/timezone"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/request"
	"github.com/mondegor/go-webcore/mrserver/request/parser"
)

const (
	moscowOffset = 3 * 60 * 60
	tokyoOffset  = 9 * 60 * 60
)

type (
	// stubLocationList - список часовых поясов с фиксированными смещениями.
	// Используется вместо timezone.LocationList, чтобы тесты не зависели
	// от наличия базы часовых поясов в системе.
	//
	// Повторяет его контракт: "Local" наружу не отдаётся (LocationByName на нём
	// промахивается), а поясом по умолчанию служит первое имя списка.
	stubLocationList struct {
		locations map[string]*time.Location
	}

	// spyLogger - логгер, считающий обращения по уровням. Нужен там, где
	// отсутствие записи в логе - часть проверяемого поведения.
	spyLogger struct {
		warns  int
		errors int
	}
)

func newStubLocationList() *stubLocationList {
	return &stubLocationList{
		// "Local" намеренно не зарегистрирован: как и timezone.LocationList,
		// стаб не выдаёт пояс процесса наружу, поэтому LocationByName на нём промахивается
		locations: map[string]*time.Location{
			"UTC":           time.UTC,
			"Europe/Moscow": time.FixedZone("Europe/Moscow", moscowOffset),
			"Asia/Tokyo":    time.FixedZone("Asia/Tokyo", tokyoOffset),
		},
	}
}

// LocationByName - повторяет контракт timezone.LocationList: при промахе
// возвращается пояс по умолчанию вместе с ошибкой.
func (l *stubLocationList) LocationByName(value string) (*time.Location, error) {
	if loc, ok := l.locations[value]; ok {
		return loc, nil
	}

	return l.Default(), errors.New("timezone not found")
}

// NameByOffset - подбирает имя пояса по смещению; летнего времени
// ни у одного пояса списка нет, поэтому isDST=true всегда даёт промах.
func (l *stubLocationList) NameByOffset(offset time.Duration, isDST bool) (name string, ok bool) {
	if isDST {
		return "", false
	}

	switch offset {
	case moscowOffset * time.Second:
		return "Europe/Moscow", true
	case tokyoOffset * time.Second:
		return "Asia/Tokyo", true
	case 0:
		return "UTC", true
	}

	return "", false
}

// Default - пояс по умолчанию, первое имя списка.
func (l *stubLocationList) Default() *time.Location {
	return l.locations["Europe/Moscow"]
}

func (l *spyLogger) Debug(_ context.Context, _ string, _ ...any)            {}
func (l *spyLogger) DebugFunc(_ context.Context, _ func() string, _ ...any) {}
func (l *spyLogger) Info(_ context.Context, _ string, _ ...any)             {}
func (l *spyLogger) Warn(_ context.Context, _ string, _ ...any)             { l.warns++ }
func (l *spyLogger) Error(_ context.Context, _ string, _ ...any)            { l.errors++ }

// Make sure the TimeZone conforms with the request.ParserTimeZone interface.
func TestTimeZoneImplementsRequestParserTimeZone(t *testing.T) {
	t.Parallel()

	assert.Implements(t, (*request.ParserTimeZone)(nil), &parser.TimeZone{})
}

// Make sure the timezone.LocationList can be used as a source of the TimeZone parser.
// Проверка компилируется только если тип go-core удовлетворяет
// интерфейсу-зависимости парсера.
func TestTimeZoneAcceptsCoreLocationList(t *testing.T) {
	t.Parallel()

	assert.NotNil(t, parser.NewTimeZone(timezone.NewLocationList(nil), mrlog.NopLogger(), "tz"))
}

func TestTimeZone_Location(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name         string
		query        string
		internalZone string
		headerSet    bool
		header       string
		want         string
		wantOffset   int
	}

	tests := []testCase{
		{
			name:       "no sources falls back to the default zone",
			want:       "Europe/Moscow",
			wantOffset: moscowOffset,
		},
		{
			name:         "internal header is used",
			internalZone: "Asia/Tokyo",
			want:         "Asia/Tokyo",
			wantOffset:   tokyoOffset,
		},
		{
			name:         "internal header takes precedence over the client header",
			internalZone: "Asia/Tokyo",
			headerSet:    true,
			header:       "UTC",
			want:         "Asia/Tokyo",
			wantOffset:   tokyoOffset,
		},
		{
			// разовое переопределение остаётся за клиентом даже на авторизованном
			// маршруте, где внутренний заголовок выставлен из профиля
			name:         "query param takes precedence over the internal header",
			query:        "?tz=UTC",
			internalZone: "Asia/Tokyo",
			headerSet:    true,
			header:       "Europe/Moscow",
			want:         "UTC",
			wantOffset:   0,
		},
		{
			// пояс профиля приложение больше не предоставляет: подбора для строгого
			// источника нет, поэтому разбор доходит до значения клиента
			name:         "unregistered internal header falls through to the client header",
			internalZone: "America/Denver",
			headerSet:    true,
			header:       "Asia/Tokyo",
			want:         "Asia/Tokyo",
			wantOffset:   tokyoOffset,
		},
		{
			name:         "unregistered internal header without a client header falls back to the default",
			internalZone: "America/Denver",
			want:         "Europe/Moscow",
			wantOffset:   moscowOffset,
		},
		{
			// внутренний заголовок несёт только имя: параметры в нём не разбираются,
			// поэтому подбора по смещению для него нет
			name:         "parameters in the internal header are not parsed",
			internalZone: "America/Denver;offset=+09:00;dst=0",
			want:         "Europe/Moscow",
			wantOffset:   moscowOffset,
		},
		{
			// имя с параметрами целиком совпадением не является даже тогда,
			// когда само имя зарегистрировано
			name:         "registered name with parameters is not a match in the internal header",
			internalZone: "Asia/Tokyo;offset=+09:00;dst=0",
			want:         "Europe/Moscow",
			wantOffset:   moscowOffset,
		},
		{
			// профиль пуст, внутренний заголовок middleware не выставил -
			// в дело вступает значение клиента
			name:       "client header applies when the internal header is absent",
			headerSet:  true,
			header:     "Asia/Tokyo",
			want:       "Asia/Tokyo",
			wantOffset: tokyoOffset,
		},
		{
			name:       "registered name from the header",
			headerSet:  true,
			header:     "Asia/Tokyo",
			want:       "Asia/Tokyo",
			wantOffset: tokyoOffset,
		},
		{
			name:       "name with parameters",
			headerSet:  true,
			header:     "Asia/Tokyo;offset=+09:00;dst=0",
			want:       "Asia/Tokyo",
			wantOffset: tokyoOffset,
		},
		{
			// имя приложению неизвестно, поэтому в дело вступает запасной путь -
			// подбор ближайшего пояса по смещению
			name:       "unregistered name falls back to the offset",
			headerSet:  true,
			header:     "America/Denver;offset=+09:00;dst=0",
			want:       "Asia/Tokyo",
			wantOffset: tokyoOffset,
		},
		{
			name:       "offset without a name",
			headerSet:  true,
			header:     "offset=+09:00;dst=0",
			want:       "Asia/Tokyo",
			wantOffset: tokyoOffset,
		},
		{
			name:       "negative offset that matches no zone falls back to the default",
			headerSet:  true,
			header:     "offset=-07:00;dst=0",
			want:       "Europe/Moscow",
			wantOffset: moscowOffset,
		},
		{
			// подбор по одному лишь смещению давал бы пояс наугад,
			// поэтому без признака летнего времени смещение не принимается
			name:       "offset without the dst flag is ignored",
			headerSet:  true,
			header:     "offset=+09:00",
			want:       "Europe/Moscow",
			wantOffset: moscowOffset,
		},
		{
			name:       "dst flag without an offset is ignored",
			headerSet:  true,
			header:     "dst=0",
			want:       "Europe/Moscow",
			wantOffset: moscowOffset,
		},
		{
			name:       "dst flag is respected",
			headerSet:  true,
			header:     "offset=+09:00;dst=1",
			want:       "Europe/Moscow",
			wantOffset: moscowOffset,
		},
		{
			name:       "offset without a sign is not accepted",
			headerSet:  true,
			header:     "offset=09:00;dst=0",
			want:       "Europe/Moscow",
			wantOffset: moscowOffset,
		},
		{
			// разбор строгий: пробел перед ключом делает негодным весь параметр,
			// поэтому запасной путь по смещению не срабатывает, а незарегистрированное
			// имя отбрасывается, и пояс сводится к значению по умолчанию
			name:       "spaces around the segments are not tolerated",
			headerSet:  true,
			header:     "America/Denver; offset=+09:00; dst=0",
			want:       "Europe/Moscow",
			wantOffset: moscowOffset,
		},
		{
			// пробел после "=" делает негодным значение параметра
			name:       "spaces around the parameter value are not tolerated",
			headerSet:  true,
			header:     "America/Denver;offset= +09:00;dst= 0",
			want:       "Europe/Moscow",
			wantOffset: moscowOffset,
		},
		{
			// имя с пробелами в список приложения не входит, поэтому пояс определяется
			// по смещению - разбор строгий и на само имя тоже
			name:       "name is not trimmed",
			headerSet:  true,
			header:     " Asia/Tokyo ;offset=+09:00;dst=0",
			want:       "Asia/Tokyo",
			wantOffset: tokyoOffset,
		},
		{
			// "Local" описывает процесс, а не клиента, и в список имён не входит
			name:       "local is not accepted from the header",
			headerSet:  true,
			header:     "Local",
			want:       "Europe/Moscow",
			wantOffset: moscowOffset,
		},
		{
			name:       "empty header falls back to the default zone",
			headerSet:  true,
			header:     "",
			want:       "Europe/Moscow",
			wantOffset: moscowOffset,
		},
		{
			name:       "utc is registered explicitly",
			headerSet:  true,
			header:     "UTC",
			want:       "UTC",
			wantOffset: 0,
		},
		{
			name:       "query param wins",
			query:      "?tz=Asia/Tokyo",
			want:       "Asia/Tokyo",
			wantOffset: tokyoOffset,
		},
		{
			name:       "query param takes precedence over the header",
			query:      "?tz=Asia/Tokyo",
			headerSet:  true,
			header:     "UTC",
			want:       "Asia/Tokyo",
			wantOffset: tokyoOffset,
		},
		{
			name:       "invalid query param is skipped in favour of the header",
			query:      "?tz=Not/AZone",
			headerSet:  true,
			header:     "Asia/Tokyo",
			want:       "Asia/Tokyo",
			wantOffset: tokyoOffset,
		},
		{
			name:       "invalid query param without a header falls back to the default",
			query:      "?tz=Not/AZone",
			want:       "Europe/Moscow",
			wantOffset: moscowOffset,
		},
		{
			// требования к параметру строже, чем к заголовку: подбора по смещению
			// у него нет, принимается только имя из списка приложения
			name:       "query param does not accept an offset",
			query:      "?tz=%2B09%3A00",
			want:       "Europe/Moscow",
			wantOffset: moscowOffset,
		},
		{
			name:       "local is not accepted from the query param",
			query:      "?tz=Local",
			want:       "Europe/Moscow",
			wantOffset: moscowOffset,
		},
		{
			// потолок - 4 сегмента, поэтому dst остаётся за ним и отбрасывается
			// вместе со смещением, а первый сегмент именем пояса не является
			name:       "segments beyond the items limit are dropped",
			headerSet:  true,
			header:     "a;b;c;offset=+09:00;dst=0",
			want:       "Europe/Moscow",
			wantOffset: moscowOffset,
		},
	}

	p := parser.NewTimeZone(newStubLocationList(), mrlog.NopLogger(), "tz")

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			r := httptest.NewRequest(http.MethodGet, "/"+tc.query, http.NoBody)
			if tc.internalZone != "" {
				r.Header.Set(mrserver.HeaderKeyInternalTimeZone, tc.internalZone)
			}

			if tc.headerSet {
				r.Header.Set(mrserver.HeaderKeyAcceptTimeZone, tc.header)
			}

			// nil-локация неотличима от UTC по имени, поэтому проверяется
			// сам объект зоны и её фактическое смещение
			loc := p.Location(r)
			require.NotNil(t, loc)
			assert.Equal(t, tc.want, loc.String())

			_, offset := time.Now().In(loc).Zone()
			assert.Equal(t, tc.wantOffset, offset)
		})
	}
}

// TestTimeZone_LocationLogging - фиксирует, что именно попадает в лог.
//
// Негодный query-параметр не логируется намеренно: годные значения клиент получает
// от самого приложения, поэтому промах - это ошибка клиента, а не сигнал серверу.
// Незарегистрированное имя в заголовке тоже молчит: клиент в поясе, которого
// приложение не предоставляет, - штатный случай. Логируется только заголовок,
// из которого не удалось извлечь ни имени, ни смещения.
func TestTimeZone_LocationLogging(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name         string
		query        string
		internalZone string
		headerSet    bool
		header       string
		wantWarns    int
	}

	tests := []testCase{
		{
			name:      "invalid query param is not logged",
			query:     "?tz=Not/AZone",
			wantWarns: 0,
		},
		{
			// расхождение профиля со списком приложения писалось бы на каждый запрос
			// такого пользователя, поэтому оно - забота проверки при старте, а не парсера
			name:         "unregistered internal header is not logged",
			internalZone: "America/Denver",
			wantWarns:    0,
		},
		{
			name:      "garbage query param is not logged",
			query:     "?tz=%D0%BC%D1%83%D1%81%D0%BE%D1%80",
			wantWarns: 0,
		},
		{
			name:      "unregistered name in the header is not logged",
			headerSet: true,
			header:    "America/Denver",
			wantWarns: 0,
		},
		{
			name:      "unusable offset in the header is not logged",
			headerSet: true,
			header:    "offset=-07:00;dst=0",
			wantWarns: 0,
		},
		{
			name:      "header without a name and an offset is logged",
			headerSet: true,
			header:    ";;;",
			wantWarns: 1,
		},
		{
			name:      "header with parameters only, none of them usable, is logged",
			headerSet: true,
			header:    "=;offset=nonsense",
			wantWarns: 1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			logger := &spyLogger{}
			p := parser.NewTimeZone(newStubLocationList(), logger, "tz")

			r := httptest.NewRequest(http.MethodGet, "/"+tc.query, http.NoBody)
			if tc.internalZone != "" {
				r.Header.Set(mrserver.HeaderKeyInternalTimeZone, tc.internalZone)
			}

			if tc.headerSet {
				r.Header.Set(mrserver.HeaderKeyAcceptTimeZone, tc.header)
			}

			p.Location(r)

			assert.Equal(t, tc.wantWarns, logger.warns)
			assert.Zero(t, logger.errors)
		})
	}
}

// TestTimeZone_BodyIsNotConsumed - фиксирует, что разбор пояса не вычитывает тело запроса.
// Парсер обязан читать только URL-запрос: r.FormValue вызвал бы ParseForm, который
// для form-encoded запроса поглощает r.Body, и обработчик получил бы пустое тело.
func TestTimeZone_BodyIsNotConsumed(t *testing.T) {
	t.Parallel()

	const body = "tz=Asia/Tokyo&payload=value"

	p := parser.NewTimeZone(newStubLocationList(), mrlog.NopLogger(), "tz")

	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// пояс из тела формы источником не является, поэтому берётся пояс по умолчанию
	assert.Equal(t, "Europe/Moscow", p.Location(r).String())

	got, err := io.ReadAll(r.Body)
	require.NoError(t, err)
	assert.Equal(t, body, string(got), "тело запроса должно остаться доступным обработчику")
}
