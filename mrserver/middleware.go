package mrserver

import (
	"net/http"
	"strings"
	"time"

	"github.com/mondegor/go-sysmess/mrlang"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrperms"
	"github.com/rs/xid"

	"github.com/mondegor/go-webcore/mrserver/mrreq"
)

// go get -u github.com/rs/xid

func MiddlewareGeneral(tr *mrlang.Translator) HttpMiddleware {
	return HttpMiddlewareFunc(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			correlationID, err := mrreq.ParseCorrelationID(r)

			if err != nil || correlationID == "" {
				correlationID = xid.New().String()
			}

			logger := mrlog.Ctx(r.Context()).With().Str("correlationID", correlationID).Logger()
			w.Header().Add(mrreq.HeaderKeyCorrelationID, correlationID)

			if err != nil {
				logger.Warn().Err(err).Msg("mrreq.ParseCorrelationID() error")
			}

			acceptLanguages := mrreq.ParseLanguage(r)
			locale := tr.FindFirstLocale(acceptLanguages...)
			logger.Debug().Str("language", locale.LangCode()).Msgf("Accept-Language: %s", strings.Join(acceptLanguages, ", "))

			srw := NewStatResponseWriter(r.Context(), w)

			defer func() {
				logger.
					Trace().
					Str("method", r.Method).
					Str("url", r.RequestURI).
					Str("userAgent", r.UserAgent()).
					Int("status", srw.statusCode).
					Int("size", srw.bytes).
					Int("elapsed_μs", int(time.Since(start).Microseconds())).
					Msg("incoming request")
			}()

			next.ServeHTTP(srw, r.WithContext(locale.WithContext(r.Context())))
		})
	})
}

func MiddlewareCheckAccess(
	section mrperms.AppSection,
	access mrperms.AccessControl,
	permission string,
	next HttpHandlerFunc,
) HttpHandlerFunc {
	privilege := section.Privilege()

	return func(w http.ResponseWriter, r *http.Request) error {
		rights := access.NewAccessRights("administrators", "guests") // :TODO: брать у пользователя

		if !rights.CheckPrivilege(privilege) && !rights.CheckPermission(permission) {
			if rights.IsGuestAccess() {
				return mrcore.FactoryErrHttpClientUnauthorized.New()
			} else {
				return mrcore.FactoryErrHttpAccessForbidden.New()
			}
		}

		return next(w, r)
	}
}
