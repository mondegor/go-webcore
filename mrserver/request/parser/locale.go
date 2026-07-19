package parser

import (
	"net/http"
	"strings"

	"github.com/mondegor/go-core/mrlocale"
	"github.com/mondegor/go-core/mrlog"
	"golang.org/x/text/language"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrserver"
)

// maxAcceptLanguageItems - предельное число элементов заголовка Accept-Language,
// разбираемых поэлементно. Клиенты перечисляют языки единицами, поэтому потолок
// на порядок выше реальных значений и на разборе годных заголовков не сказывается.
// Нужен из-за того, что поэлементный разбор включается негодным заголовком, который
// задаёт клиент: без потолка заголовок вида "zz,zz,zz,..." размером с лимит
// net/http на заголовки давал бы сотни тысяч разборов на один запрос.
const (
	maxAcceptLanguageItems = 8
)

type (
	// Locale - определяет локаль и язык пользователя из HTTP-запроса.
	Locale struct {
		pool          *mrlocale.Pool
		logger        mrlog.Logger
		paramNameLang string
	}
)

// NewLocale - создаёт объект Locale.
func NewLocale(
	pool *mrlocale.Pool,
	logger mrlog.Logger,
	paramNameLang string,
) *Locale {
	return &Locale{
		pool:          pool,
		logger:        logger,
		paramNameLang: paramNameLang,
	}
}

// Language - возвращает код языка из HTTP запроса.
func (p *Locale) Language(r *http.Request) string {
	return p.locale(r).Language()
}

// Localizer - возвращает локализатор для HTTP запроса.
func (p *Locale) Localizer(r *http.Request) mrcore.Localizer {
	return p.locale(r)
}

// locale - определяет локаль запроса по приоритету источников:
//  1. query-параметр (?lang) - одиночный тег; при валидности добавляется первым в список предпочтений;
//  2. заголовок Accept-Language - список тегов, отсортированный по убыванию веса q; для авторизованных
//     запросов middleware подставляет сюда язык из профиля пользователя (клиентский заголовок игнорируется).
//
// Невалидное значение источника логируется и пропускается (не прерывает разбор), а из заголовка
// при этом извлекаются языки тех его элементов, которые разобрать удалось. Итоговый список
// предпочтений передаётся в language.Matcher пула поддерживаемых языков, который выбирает наиболее
// близкий поддерживаемый язык. Теги с регионом (напр. fr-CH, en-US) сводятся к базовому языку.
// При пустом списке предпочтений, а также когда ни один тег не подошёл, возвращается язык по умолчанию.
func (p *Locale) locale(r *http.Request) *mrlocale.Localizer {
	langs := make([]language.Tag, 0, 2)

	if langCode := r.URL.Query().Get(p.paramNameLang); langCode != "" {
		if lang, err := language.Parse(langCode); err != nil {
			p.logger.Warn(
				r.Context(),
				"Language param is incorrect",
				"param", p.paramNameLang,
				"lang", langCode,
			)
		} else {
			langs = append(langs, lang)
		}
	}

	if acceptLanguage := r.Header.Get(mrserver.HeaderKeyAcceptLanguage); acceptLanguage != "" {
		lang, _, err := language.ParseAcceptLanguage(acceptLanguage)
		if err != nil {
			// в заголовке содержится сломанный элемент, поэтому языки,
			// перечисленные рядом с этим элементом, восстанавливаются отдельно
			lang = parseAcceptLanguageByItem(acceptLanguage)

			p.logger.Warn(
				r.Context(),
				"Header is incorrect",
				"header", mrserver.HeaderKeyAcceptLanguage,
				"lang", acceptLanguage,
				"languages", lang,
			)
		}

		langs = append(langs, lang...)
	}

	p.logger.Debug(r.Context(), "Parse locale", "languages", langs)

	return p.pool.Localizer(langs...)
}

// parseAcceptLanguageByItem - разбирает заголовок Accept-Language поэлементно, пропуская
// элементы, которые разобрать не удалось, и возвращает языки уцелевших элементов
// в порядке убывания веса q. Если не уцелел ни один элемент, возвращает пустой список.
//
// Применяется, когда заголовок не удалось разобрать целиком: language.ParseAcceptLanguage
// на любом неизвестном языке или негодном весе возвращает пустой список, поэтому один
// негодный элемент лишал бы клиента и всех остальных перечисленных им языков.
//
// Рассматриваются только первые maxAcceptLanguageItems элементов, остальные отбрасываются
// без разбора.
func parseAcceptLanguageByItem(acceptLanguage string) []language.Tag {
	// SplitN с запасом в один элемент: хвост за потолком попадает в последний
	// элемент целиком и отбрасывается вместе с ним, не порождая разборов по числу запятых
	items := strings.SplitN(acceptLanguage, ",", maxAcceptLanguageItems+1)

	if len(items) > maxAcceptLanguageItems {
		items = items[:maxAcceptLanguageItems]
	}

	survived := make([]string, 0, len(items))

	for _, item := range items {
		// элемент проверяется тем же разбором, что и весь заголовок, поэтому
		// негодный вес отсеивает элемент наравне с неизвестным языком
		if tags, _, err := language.ParseAcceptLanguage(item); err != nil || len(tags) == 0 {
			continue
		}

		survived = append(survived, item)
	}

	// порядок по весу восстанавливается общим разбором уцелевших элементов;
	// каждый из них уже разобран по отдельности, поэтому ошибка здесь означает
	// лишь то, что языков не осталось, и обрабатывается пустым списком
	tags, _, _ := language.ParseAcceptLanguage(strings.Join(survived, ","))

	return tags
}
