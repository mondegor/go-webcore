package mrserver

import (
    "net"
    "net/http"

    "github.com/mondegor/go-sysmess/mrlang"
    "github.com/mondegor/go-webcore/mrcore"
    "github.com/mondegor/go-webcore/mrctx"
    "github.com/mondegor/go-webcore/mrreq"
)

func MiddlewareFirst(l mrcore.Logger) mrcore.HttpMiddleware {
    return mrcore.HttpMiddlewareFunc(func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            correlationId, err := mrreq.CorrelationId(r)

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

func MiddlewareUserIp() mrcore.HttpMiddleware {
    return mrcore.HttpMiddlewareFunc(func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            logger := mrctx.Logger(r.Context())
            logger.Debug("Exec MiddlewareUserIp")

            userIp, err := mrreq.UserIp(r)

            if err != nil {
                logger.Warn(err.Error())
                userIp = net.IP{}
            }

            logger.Info("UserIp: %s", userIp.String())

            ctx := mrctx.WithUserIp(r.Context(), userIp)

            next.ServeHTTP(w, r.WithContext(ctx))
        })
    })
}

func MiddlewareAcceptLanguage(translator *mrlang.Translator) mrcore.HttpMiddleware {
    return mrcore.HttpMiddlewareFunc(func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            logger := mrctx.Logger(r.Context())
            logger.Debug("Exec MiddlewareAcceptLanguage")

            acceptLanguages := mrreq.Language(r)
            locale := translator.FindFirstLocale(acceptLanguages...)

            logger.Info("Accept-Language: %v; Set-Language: %s", acceptLanguages, locale.LangCode())
            ctx := mrctx.WithLocale(r.Context(), locale)

            next.ServeHTTP(w, r.WithContext(ctx))
        })
    })
}

func MiddlewarePlatform(defaultPlatform string) mrcore.HttpMiddleware {
    return mrcore.HttpMiddlewareFunc(func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            logger := mrctx.Logger(r.Context())
            logger.Debug("Exec MiddlewarePlatform")

            platform, err := mrreq.Platform(r)

            if err != nil {
                logger.Warn(err.Error())
            }

            if platform == "" {
                platform = defaultPlatform
            }

            logger.Info("Platform: %s", platform)

            ctx := mrctx.WithPlatform(r.Context(), platform)

            next.ServeHTTP(w, r.WithContext(ctx))
        })
    })
}

func MiddlewareAuthenticateUser() mrcore.HttpMiddleware {
    return mrcore.HttpMiddlewareFunc(func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            logger := mrctx.Logger(r.Context())
            logger.Debug("Exec MiddlewareAuthenticateUser")

            next.ServeHTTP(w, r)
        })
    })
}
