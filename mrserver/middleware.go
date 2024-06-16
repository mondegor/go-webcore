package mrserver

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/mondegor/go-sysmess/mrlang"
	"github.com/rs/xid"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mridempotency"
	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrperms"

	"github.com/mondegor/go-webcore/mrserver/mrreq"
)

// go get -u github.com/rs/xid

// MiddlewareGeneral - comment func.
func MiddlewareGeneral(
	tr *mrlang.Translator,
	statFunc func(l mrlog.Logger, start time.Time, sr *StatRequest, sw *StatResponseWriter),
) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
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

			r = r.WithContext(locale.WithContext(logger.WithContext(r.Context())))

			sr := NewStatRequest(r)
			sw := NewStatResponseWriter(
				w,
				func(buf []byte) {
					logger.Trace().Bytes("response", buf).Msg("write response")
				},
			)

			defer func() {
				statFunc(logger, start, sr, sw)
			}()

			next.ServeHTTP(sw, r)
		})
	}
}

// MiddlewareRecoverHandler - comment func.
func MiddlewareRecoverHandler(isDebug bool, fatalFunc http.HandlerFunc) func(next http.Handler) http.Handler {
	if fatalFunc == nil {
		fatalFunc = func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func(ctx context.Context) {
				if rvr := recover(); rvr != nil {
					if rvr == http.ErrAbortHandler { //nolint:errorlint
						// we don't recover http.ErrAbortHandler so the response
						// to the client is aborted, this should not be logged
						panic(rvr)
					}

					if isDebug {
						os.Stderr.Write([]byte(fmt.Sprintf("%+v", r)))
						os.Stderr.Write(debug.Stack())
					} else {
						mrlog.Ctx(ctx).Error().
							Str("panic", fmt.Sprintf("%+v", r)).
							Bytes("CallStack", debug.Stack()).
							Send()
					}

					fatalFunc(w, r)
				}
			}(r.Context())

			next.ServeHTTP(w, r)
		})
	}
}

// MiddlewareHandlerAdapter - comment func.
func MiddlewareHandlerAdapter(s ErrorResponseSender) func(next HttpHandlerFunc) http.HandlerFunc {
	return func(next HttpHandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if err := next(w, r); err != nil {
				s.SendError(w, r, err)
			}
		}
	}
}

// MiddlewareHandlerCheckAccess - comment func.
func MiddlewareHandlerCheckAccess(
	handlerName string,
	access mrperms.AccessRightsFactory,
	privilege, permission string,
) func(next HttpHandlerFunc) HttpHandlerFunc {
	return func(next HttpHandlerFunc) HttpHandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) error {
			mrlog.Ctx(r.Context()).Debug().Str("handler", handlerName).Msg("exec handler")

			rights := access.NewAccessRights("administrators", "guests") // TODO: брать у пользователя

			if rights.CheckPrivilege(privilege) && rights.CheckPermission(permission) {
				return next(w, r)
			}

			if rights.IsGuestAccess() {
				return mrcore.ErrHttpClientUnauthorized.New()
			}

			return mrcore.ErrHttpAccessForbidden.New()
		}
	}
}

// MiddlewareHandlerIdempotency - comment func.
func MiddlewareHandlerIdempotency(provider mridempotency.Provider, sender ResponseSender) func(next HttpHandlerFunc) HttpHandlerFunc {
	return func(next HttpHandlerFunc) HttpHandlerFunc {
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

			if err := next(sw, r); err != nil {
				return err
			}

			if err := provider.Store(r.Context(), idempotencyKey, sw); err != nil {
				mrlog.Ctx(r.Context()).Error().Err(err).Send()
			}

			return nil
		}
	}
}
