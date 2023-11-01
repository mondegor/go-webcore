package mrserver

import (
    "net/http"

    "github.com/mondegor/go-sysmess/mrlang"
    "github.com/mondegor/go-webcore/mrcore"
    "github.com/mondegor/go-webcore/mrctx"
    "github.com/mondegor/go-webcore/mrreq"
)

func MiddlewareFirst(l mrcore.Logger) mrcore.HttpMiddleware {
    return mrcore.HttpMiddlewareFunc(func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            correlationId, err := mrreq.ParseCorrelationId(r)

            if err != nil {
                l.Warn(err.Error())
            }

            if correlationId == "" {
                correlationId = mrctx.GenCorrelationId()
            }

            logger := l.With(correlationId)

            logger.Debug("Exec MiddlewareFirst")
            logger.Info("CorrelationID: %s", correlationId)

            ctx := mrctx.WithCorrelationId(r.Context(), correlationId)
            ctx = mrctx.WithLogger(ctx, logger)

            next.ServeHTTP(w, r.WithContext(ctx))
        })
    })
}

func MiddlewareAcceptLanguage(translator *mrlang.Translator) mrcore.HttpMiddleware {
    return mrcore.HttpMiddlewareFunc(func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            logger := mrctx.Logger(r.Context())
            logger.Debug("Exec MiddlewareAcceptLanguage")

            acceptLanguages := mrreq.ParseLanguage(r)
            locale := translator.FindFirstLocale(acceptLanguages...)

            logger.Info("Accept-Language: %v; Set-Language: %s", acceptLanguages, locale.LangCode())
            ctx := mrctx.WithLocale(r.Context(), locale)

            next.ServeHTTP(w, r.WithContext(ctx))
        })
    })
}
