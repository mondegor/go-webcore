package mrparser

import (
	"fmt"
	"net/http"

	"github.com/mondegor/go-sysmess/mrlocale"
	"github.com/mondegor/go-sysmess/mrlog"
	"golang.org/x/text/language"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrserver/mrreq"
)

type (
	// Locale - comment struct.
	Locale struct {
		pool          *mrlocale.Pool
		logger        mrlog.Logger
		paramNameLang string
	}
)

// NewLocale - создаёт объект Locale.
func NewLocale(pool *mrlocale.Pool, logger mrlog.Logger, paramNameLang string) *Locale {
	return &Locale{
		pool:          pool,
		logger:        logger,
		paramNameLang: paramNameLang,
	}
}

// Language - comment method.
func (p *Locale) Language(r *http.Request) string {
	return p.locale(r).Language()
}

// Localizer - comment method.
func (p *Locale) Localizer(r *http.Request) mrcore.Localizer {
	return p.locale(r)
}

func (p *Locale) locale(r *http.Request) *mrlocale.Localizer {
	langs := make([]language.Tag, 0, 2)

	if langCode := r.FormValue(p.paramNameLang); langCode != "" {
		if lang, err := language.Parse(langCode); err != nil {
			p.logger.Warn(r.Context(), fmt.Sprintf("Language param %s with value %s is incorrect", p.paramNameLang, langCode))
		} else {
			langs = append(langs, lang)
		}
	}

	if acceptLanguage := r.Header.Get(mrreq.HeaderKeyAcceptLanguage); acceptLanguage != "" {
		if lang, _, err := language.ParseAcceptLanguage(acceptLanguage); err != nil {
			p.logger.Warn(r.Context(), fmt.Sprintf("Header %s with value %s is incorrect", mrreq.HeaderKeyAcceptLanguage, acceptLanguage))
		} else {
			langs = append(langs, lang...)
		}
	}

	p.logger.Debug(r.Context(), "Parse locale", "languages", langs)

	return p.pool.Localizer(langs...)
}
