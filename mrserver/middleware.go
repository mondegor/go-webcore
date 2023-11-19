package mrserver

import (
	"net/http"

	"github.com/mondegor/go-sysmess/mrlang"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrcrypto"
	"github.com/mondegor/go-webcore/mrctx"
	"github.com/mondegor/go-webcore/mrreq"
)

func MiddlewareFirst(l mrcore.Logger, t *mrlang.Translator, v mrcore.Validator) mrcore.HttpMiddleware {
	return mrcore.HttpMiddlewareFunc(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			correlationID, err := mrreq.ParseCorrelationID(r)

			if err != nil || correlationID == "" {
				correlationID = mrcrypto.GenTokenHexWithDelimiter(9, 4)
			}

			logger := l.With(correlationID)

			if err != nil {
				logger.DisableFileLine().Warn(err)
			}

			logger.Debug("Exec MiddlewareFirst")
			logger.Debug("%s %s", r.Method, r.RequestURI)
			logger.Debug("CorrelationID: %s", correlationID)

			acceptLanguages := mrreq.ParseLanguage(r)
			locale := t.FindFirstLocale(acceptLanguages...)

			logger.Debug("Accept-Language: %v; Set-Language: %s", acceptLanguages, locale.LangCode())

			ctx := mrctx.WithClientTools(r.Context(), correlationID, logger, locale, v)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})
}
