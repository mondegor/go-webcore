package mrserver

import (
	"context"
	"errors"
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

// :TODO: вынести в настройки

const (
	traceRequestBodyMaxLen  = 2048
	traceResponseBodyMaxLen = 2048
)

// go get -u github.com/rs/xid

// MiddlewareGeneral - промежуточный обработчик, который устанавливает в контекст
// requestId, language, logger. А также другие параметры, которые используются в статистике запросов.
func MiddlewareGeneral(
	tr *mrlang.Translator,
	observeRequestFunc func(l mrlog.Logger, start time.Time, sr *StatRequestReader, sw *StatResponseWriter),
) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			requestID := xid.New().String()
			correlationID, err := mrreq.ParseCorrelationID(r.Header)

			if err == nil && correlationID != "" {
				requestID += "-" + correlationID
			}

			logger := mrlog.Ctx(r.Context()).With().Str("requestId", requestID).Logger()

			if err != nil {
				logger.Warn().Err(err).Msg("mrreq.ParseCorrelationID()")
			}

			w.Header().Add(mrreq.HeaderKeyRequestID, requestID)

			acceptLanguages := mrreq.ParseLanguage(r.Header)
			logger.Debug().Msgf("Accept-Language: %s", strings.Join(acceptLanguages, ", "))

			locale := tr.FindFirstLocale(acceptLanguages...)
			logger.Info().
				Str("method", r.Method).
				Str("uri", r.RequestURI).
				Str("language", locale.LangCode()).
				Msg("request")

			r = r.WithContext(locale.WithContext(logger.WithContext(r.Context())))
			sr := NewStatRequestReader(r, traceRequestBodyMaxLen)
			sw := NewStatResponseWriter(w, traceResponseBodyMaxLen)

			defer observeRequestFunc(logger, start, sr, sw)

			next.ServeHTTP(sw, sr.Request())
		})
	}
}

// MiddlewareRecoverHandler - промежуточный обработчик для перехвата panic.
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

					errorMessage := fmt.Sprintf("%s method %s %s; panic: %v", r.Proto, r.Method, r.URL, rvr)

					if isDebug {
						os.Stderr.Write([]byte(errorMessage + "\n"))
						os.Stderr.Write(debug.Stack())
					} else {
						mrlog.Ctx(ctx).
							Error().
							Err(errors.New(errorMessage)).
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

// MiddlewareHandlerAdapter - переходник с HttpHandlerFunc на http.HandlerFunc.
func MiddlewareHandlerAdapter(s ErrorResponseSender) func(next HttpHandlerFunc) http.HandlerFunc {
	return func(next HttpHandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if err := next(w, r); err != nil {
				s.SendError(w, r, err)
			}
		}
	}
}

// MiddlewareHandlerCheckAccess - промежуточный обработчик проверки доступа к секции и конечному обработчику.
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

// MiddlewareHandlerIdempotency - промежуточный обработчик для организации идемпотентных запросов.
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

			cachedResponse, err := provider.Get(r.Context(), idempotencyKey)
			if err != nil {
				return err
			}

			if cachedResponse != nil {
				return sender.SendBytes(
					w,
					cachedResponse.StatusCode(),
					cachedResponse.Content(),
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
}
