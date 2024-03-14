package mrserver

import (
	"net/http"
	"strings"
	"time"

	"github.com/mondegor/go-sysmess/mrlang"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mridempotency"
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
				logger.Warn().Err(err).Msg("mrreq.ParseCorrelationID()")
			}

			acceptLanguages := mrreq.ParseLanguage(r)
			locale := tr.FindFirstLocale(acceptLanguages...)
			logger.Debug().
				Str("language", locale.LangCode()).
				Msgf("Accept-Language: %s", strings.Join(acceptLanguages, ", "))

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

func MiddlewareIdempotency(provider mridempotency.Provider, sender ResponseSender, next HttpHandlerFunc) HttpHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		idempotencyKey := r.Header.Get(mrreq.HeaderKeyIdempotencyKey)

		if idempotencyKey == "" {
			return next(w, r)
		}

		if err := provider.Validate(idempotencyKey); err != nil {
			return err
		}

		if cachedResponse, err := provider.Get(r.Context(), idempotencyKey); err != nil {
			return err
		} else if cachedResponse != nil {
			return sender.SendBytes(
				w,
				cachedResponse.StatusCode(),
				cachedResponse.Body(),
			)
		}

		unlock, err := provider.Lock(r.Context(), idempotencyKey)

		if err != nil {
			return err
		}

		defer unlock()

		sw := NewCacheableResponseWriter(w)

		if err = next(sw, r); err != nil {
			return err
		}

		if err = provider.Store(r.Context(), idempotencyKey, sw); err != nil {
			mrlog.Ctx(r.Context()).Error().Err(err).Send()
		}

		return nil
	}
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
