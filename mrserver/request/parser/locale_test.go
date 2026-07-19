package parser_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mondegor/go-core/mrlocale"
	"github.com/mondegor/go-core/mrlog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"

	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/request"
	"github.com/mondegor/go-webcore/mrserver/request/parser"
)

type (
	// stubMessageProvider - провайдер локализации без переводов.
	// Тестам парсера важен только выбор языка, а не содержимое сообщений.
	stubMessageProvider struct{}
)

func (p stubMessageProvider) Domains() []string {
	return []string{mrlocale.DefaultMessagesDomain, mrlocale.DefaultErrorsDomain}
}

func (p stubMessageProvider) Localize(_ string, _ language.Tag, msg string, _ []any) string {
	return msg
}

// newTestLocalePool - создаёт пул с языками ru (по умолчанию), en и fr.
func newTestLocalePool(t *testing.T) *mrlocale.Pool {
	t.Helper()

	bundle, err := mrlocale.NewBundle(
		[]string{"ru", "en", "fr"},
		mrlocale.WithMessageProvider(func(_ []language.Tag) (mrlocale.MessageProvider, error) {
			return stubMessageProvider{}, nil
		}),
	)
	require.NoError(t, err)

	return mrlocale.NewPool(bundle)
}

// Make sure the Locale conforms with the request.ParserLocale interface.
func TestLocaleImplementsRequestParserLocale(t *testing.T) {
	t.Parallel()

	assert.Implements(t, (*request.ParserLocale)(nil), &parser.Locale{})
}

func TestLocale_Language(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name           string
		query          string
		acceptLanguage string
		want           string
	}

	tests := []testCase{
		{
			name: "no sources falls back to the default language",
			want: "ru",
		},
		{
			name:  "query param wins",
			query: "?lang=en",
			want:  "en",
		},
		{
			name:           "query param takes precedence over the header",
			query:          "?lang=en",
			acceptLanguage: "fr",
			want:           "en",
		},
		{
			name:           "header is used when the query param is absent",
			acceptLanguage: "fr",
			want:           "fr",
		},
		{
			name:           "header respects q-weights",
			acceptLanguage: "fr;q=0.8,en;q=0.9",
			want:           "en",
		},
		{
			// тег с регионом доходит до базового языка: language.Matcher находит fr,
			// но возвращает тег с расширением региона ("fr-u-rg-chzzzz"), поэтому
			// mrlocale.Pool.Localizer выбирает язык по индексу совпадения, а не по тегу
			name:  "regional tag resolves to the base language",
			query: "?lang=fr-CH",
			want:  "fr",
		},
		{
			name:           "invalid query param is skipped in favour of the header",
			query:          "?lang=not-a-language-tag",
			acceptLanguage: "en",
			want:           "en",
		},
		{
			// подходящего языка в пуле нет: language.Matcher сообщает confidence=No,
			// такой промах отсекается отдельно, поэтому выдаётся язык по умолчанию,
			// а не первый язык списка, на который указывает индекс матчера
			name:           "unsupported language falls back to the default",
			acceptLanguage: "de",
			want:           "ru",
		},
	}

	p := parser.NewLocale(newTestLocalePool(t), mrlog.NopLogger(), "lang")

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			r := httptest.NewRequest(http.MethodGet, "/"+tc.query, http.NoBody)
			if tc.acceptLanguage != "" {
				r.Header.Set(mrserver.HeaderKeyAcceptLanguage, tc.acceptLanguage)
			}

			assert.Equal(t, tc.want, p.Language(r))
		})
	}
}

// TestLocale_LanguagePartiallyIncorrectHeader - фиксирует, что негодный элемент заголовка
// Accept-Language не отбрасывает языки, перечисленные рядом с ним.
//
// language.ParseAcceptLanguage разбирает заголовок по принципу "всё или ничего": на любом
// неизвестном языке или негодном весе он возвращает пустой список, поэтому клиент, явно
// запросивший поддерживаемый язык, получал бы язык по умолчанию.
func TestLocale_LanguagePartiallyIncorrectHeader(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name           string
		acceptLanguage string
		want           string
	}

	tests := []testCase{
		{
			// регион EN не существует, поэтому весь элемент отбрасывается
			name:           "unknown region does not discard the next language",
			acceptLanguage: "en-EN,fr;q=0.9",
			want:           "fr",
		},
		{
			name:           "unknown language does not discard the next language",
			acceptLanguage: "zz,fr;q=0.9",
			want:           "fr",
		},
		{
			name:           "invalid weight does not discard the next language",
			acceptLanguage: "en;q=abc,fr;q=0.9",
			want:           "fr",
		},
		{
			// вес уцелевших элементов учитывается, как и при обычном разборе
			name:           "weights of the survived languages are respected",
			acceptLanguage: "zz,fr;q=0.8,en;q=0.9",
			want:           "en",
		},
		{
			name:           "fully incorrect header falls back to the default",
			acceptLanguage: "garbage!!!",
			want:           "ru",
		},
		{
			// поэлементно разбираются только первые 8 элементов, поэтому en,
			// стоящий за потолком, до пула не доходит
			name:           "languages beyond the items limit are dropped",
			acceptLanguage: strings.Repeat("zz,", 8) + "en",
			want:           "ru",
		},
		{
			// последний элемент в пределах потолка ещё учитывается
			name:           "the last language within the items limit survives",
			acceptLanguage: strings.Repeat("zz,", 7) + "en",
			want:           "en",
		},
	}

	p := parser.NewLocale(newTestLocalePool(t), mrlog.NopLogger(), "lang")

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			r := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
			r.Header.Set(mrserver.HeaderKeyAcceptLanguage, tc.acceptLanguage)

			assert.Equal(t, tc.want, p.Language(r))
		})
	}
}

// TestLocale_BodyIsNotConsumed - фиксирует, что разбор локали не вычитывает тело запроса.
// Парсер обязан читать только URL-запрос: r.FormValue вызвал бы ParseForm, который
// для form-encoded запроса поглощает r.Body, и обработчик получил бы пустое тело.
func TestLocale_BodyIsNotConsumed(t *testing.T) {
	t.Parallel()

	const body = "lang=en&payload=value"

	p := parser.NewLocale(newTestLocalePool(t), mrlog.NopLogger(), "lang")

	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// язык из тела формы источником не является, поэтому берётся язык по умолчанию
	assert.Equal(t, "ru", p.Language(r))

	got, err := io.ReadAll(r.Body)
	require.NoError(t, err)
	assert.Equal(t, body, string(got), "тело запроса должно остаться доступным обработчику")
}
