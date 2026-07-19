package parser_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
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

type (
	// stubLocationList - список часовых поясов с фиксированными смещениями.
	// Используется вместо timezone.LocationList, чтобы тесты не зависели
	// от наличия базы часовых поясов в системе.
	stubLocationList struct {
		locations map[string]*time.Location
	}
)

func newStubLocationList() *stubLocationList {
	return &stubLocationList{
		locations: map[string]*time.Location{
			// UTC и Local повторяют timezone.NewLocationList, который регистрирует
			// оба имени всегда, независимо от переданного списка
			"UTC":           time.UTC,
			"Local":         time.Local,
			"Europe/Moscow": time.FixedZone("Europe/Moscow", 3*60*60),
			"Asia/Tokyo":    time.FixedZone("Asia/Tokyo", 9*60*60),
		},
	}
}

// LocationByName - повторяет контракт timezone.LocationList: при промахе
// возвращается пояс по умолчанию вместе с ошибкой.
func (l *stubLocationList) LocationByName(value string) (*time.Location, error) {
	if loc, ok := l.locations[value]; ok {
		return loc, nil
	}

	return time.UTC, errors.New("timezone not found")
}

// Make sure the User conforms with the request.ParserUser interface.
func TestUserImplementsRequestParserUser(t *testing.T) {
	t.Parallel()

	assert.Implements(t, (*request.ParserUser)(nil), &parser.User{})
}

// Make sure the timezone.LocationList can be used as a source of the User parser.
// Проверка компилируется только если тип go-core удовлетворяет
// интерфейсу-зависимости парсера.
func TestUserAcceptsCoreLocationList(t *testing.T) {
	t.Parallel()

	assert.NotNil(t, parser.NewUser(&timezone.LocationList{}, mrlog.NopLogger()))
}

func TestUser_SessionID(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name      string
		headerSet bool
		header    string
		want      string
	}

	tests := []testCase{
		{
			name:      "header is set",
			headerSet: true,
			header:    "8f14e45f",
			want:      "8f14e45f",
		},
		{
			name:      "header is empty",
			headerSet: true,
			header:    "",
			want:      "",
		},
		{
			name:      "header is absent",
			headerSet: false,
			want:      "",
		},
	}

	p := parser.NewUser(newStubLocationList(), mrlog.NopLogger())

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			r := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
			if tc.headerSet {
				r.Header.Set(mrserver.HeaderKeySessionID, tc.header)
			}

			assert.Equal(t, tc.want, p.SessionID(r))
		})
	}
}

// TestUser_LocationWithoutLocationList - проверяет, что парсер, созданный без списка
// часовых поясов, не паникует, а молча сводит любой пояс к time.UTC.
func TestUser_LocationWithoutLocationList(t *testing.T) {
	t.Parallel()

	p := parser.NewUser(nil, mrlog.NopLogger())

	r := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	r.Header.Set(mrserver.HeaderKeyTimeZone, "Europe/Moscow")

	assert.Equal(t, time.UTC, p.Location(r))
}

// TestUser_LocationLocalIsRegistered - фиксирует, что служебное имя "Local" обычным
// незарегистрированным именем не является: timezone.NewLocationList регистрирует его
// всегда, поэтому парсер отдаёт часовой пояс процесса, а не time.UTC.
//
// Сам парсер имя не разбирает и служебные имена не отсеивает - решение целиком
// за locationList, поэтому "Local" доходит до результата наравне с IANA-именем.
func TestUser_LocationLocalIsRegistered(t *testing.T) {
	t.Parallel()

	p := parser.NewUser(newStubLocationList(), mrlog.NopLogger())

	r := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	r.Header.Set(mrserver.HeaderKeyTimeZone, "Local")

	assert.Equal(t, time.Local, p.Location(r))
}

func TestUser_Location(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name       string
		headerSet  bool
		header     string
		want       string
		wantOffset int // смещение зоны в секундах
	}

	tests := []testCase{
		{
			name:       "registered zone",
			headerSet:  true,
			header:     "Europe/Moscow",
			want:       "Europe/Moscow",
			wantOffset: 3 * 60 * 60,
		},
		{
			name:       "another registered zone",
			headerSet:  true,
			header:     "Asia/Tokyo",
			want:       "Asia/Tokyo",
			wantOffset: 9 * 60 * 60,
		},
		{
			name:       "UTC is registered explicitly",
			headerSet:  true,
			header:     "UTC",
			want:       "UTC",
			wantOffset: 0,
		},
		{
			name:       "unregistered zone falls back to UTC",
			headerSet:  true,
			header:     "Not/AZone",
			want:       "UTC",
			wantOffset: 0,
		},
		{
			name:       "header is empty",
			headerSet:  true,
			header:     "",
			want:       "UTC",
			wantOffset: 0,
		},
		{
			name:       "header is absent",
			headerSet:  false,
			want:       "UTC",
			wantOffset: 0,
		},
	}

	p := parser.NewUser(newStubLocationList(), mrlog.NopLogger())

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			r := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
			if tc.headerSet {
				r.Header.Set(mrserver.HeaderKeyTimeZone, tc.header)
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
