package mrserver

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"runtime/debug"
	"time"

	"github.com/mondegor/go-sysmess/mrerr/mr"
	"github.com/mondegor/go-sysmess/mrlog"

	"github.com/mondegor/go-webcore/mraccess"
	"github.com/mondegor/go-webcore/mridempotency"
	"github.com/mondegor/go-webcore/mrserver/mrreq"
	"github.com/mondegor/go-webcore/mrserver/observe"
)

const (
	// PrivilegePublic - привилегия для всех.
	PrivilegePublic = "public"

	// PermissionAnyUser - разрешение для любого пользователя.
	PermissionAnyUser = "any-user"

	// PermissionGuestOnly - разрешение только для гостя.
	PermissionGuestOnly = "guest-only"
)

const (
	// :TODO: вынести в настройки.
	traceRequestBodyMaxLen  = 2048
	traceResponseBodyMaxLen = 2048
)

type (
	traceManager interface {
		WithCorrelationID(ctx context.Context, id string) context.Context
		WithGeneratedRequestID(ctx context.Context) context.Context
		RequestID(ctx context.Context) string
	}
)

// go get -u github.com/rs/xid

// MiddlewareRecoverHandler - промежуточный обработчик для перехвата panic.
func MiddlewareRecoverHandler(logger mrlog.Logger, isDebug bool, fatalFunc http.HandlerFunc) func(next http.Handler) http.Handler {
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

					errorMessage := fmt.Sprintf("%s method %s %s", r.Proto, r.Method, r.URL)

					if isDebug {
						os.Stderr.Write([]byte(errorMessage))
						os.Stderr.Write([]byte(fmt.Sprintf("; panic: %v\n", rvr)))
						os.Stderr.Write(debug.Stack())
					} else {
						logger.Error(
							ctx,
							"MiddlewareRecoverHandler",
							"error",
							mr.ErrInternalCaughtPanic.New(
								errorMessage,
								rvr,
								string(debug.Stack()),
							),
						)
					}

					fatalFunc(w, r)
				}
			}(r.Context())

			next.ServeHTTP(w, r)
		})
	}
}

// MiddlewareRequestID - промежуточный обработчик,
// который устанавливает в контекст requestId, correlationId.
func MiddlewareRequestID(logger mrlog.Logger, traceManager traceManager) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := traceManager.WithGeneratedRequestID(r.Context())
			w.Header().Set(mrreq.HeaderKeyRequestID, traceManager.RequestID(ctx))

			if correlationID, err := mrreq.ParseCorrelationID(r.Header); err != nil {
				logger.Warn(ctx, "MiddlewareRequestID", "error", err)
			} else if correlationID != "" {
				ctx = traceManager.WithCorrelationID(ctx, correlationID)
				w.Header().Set(mrreq.HeaderKeyCorrelationID, correlationID)
			}

			// необходимо гарантировать, чтобы этот заголовок не был передан из вне
			r.Header.Del(mrreq.HeaderKeyUserIDSlashGroup)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// MiddlewareObserver - промежуточный обработчик, который собирает статистику запросов.
func MiddlewareObserver(
	logger mrlog.Logger,
	observer RequestStat,
) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			logger.Info(r.Context(), "REQUEST",
				"method", r.Method,
				"uri", r.RequestURI,
			)

			sr := observe.NewRequestReader(r, traceRequestBodyMaxLen)
			sw := observe.NewResponseWriter(w, traceResponseBodyMaxLen)

			defer func() {
				observer.Emit(
					sr.Request(),
					sr.Content(),
					sr.Size(),
					sw.Content(),
					sw.Size(),
					time.Since(start),
					sw.StatusCode(),
				)
			}()

			next.ServeHTTP(sw, sr.Request())
		})
	}
}

// MiddlewareHandlerAdapter - переходник с HttpHandlerFunc на http.HandlerFunc.
func MiddlewareHandlerAdapter(errSender ErrorResponseSender) func(next HttpHandlerFunc) http.HandlerFunc {
	return func(next HttpHandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if err := next(w, r); err != nil {
				if mr.ErrUseCaseEntityNotFound.Is(err) { // подменяются только необёрнутые ошибки этого типа
					err = mr.ErrHttpResourceNotFound.New()
				}

				errSender.SendError(w, r, err)
			}
		}
	}
}

// MiddlewareHandlerCheckAccessToken - промежуточный обработчик запрещает доступ к обработчику авторизованному пользователю.
func MiddlewareHandlerCheckAccessToken(logger mrlog.Logger, handlerName string) func(next HttpHandlerFunc) HttpHandlerFunc {
	return func(next HttpHandlerFunc) HttpHandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) error {
			logger.Debug(r.Context(), "MiddlewareHandlerCheckAccessToken", "handler", handlerName)

			if accessToken := mrreq.ParseAccessToken(r.Header); accessToken != "" {
				return mr.ErrHttpAccessForbidden.New()
			}

			return next(w, r)
		}
	}
}

// MiddlewareHandlerCheckAccess - промежуточный обработчик проверки доступа к секции и конечному обработчику.
func MiddlewareHandlerCheckAccess(
	logger mrlog.Logger,
	handlerName, privilege, permission string,
	userProvider mraccess.MemberProvider,
	userGroups mraccess.RightsGetter,
) func(next HttpHandlerFunc) HttpHandlerFunc {
	return func(next HttpHandlerFunc) HttpHandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) error {
			ctx := r.Context()
			logger.Debug(
				ctx,
				"MiddlewareHandlerCheckAccess",
				"handler", handlerName,
				"privilege", privilege,
				"permission", permission,
			)

			accessToken := mrreq.ParseAccessToken(r.Header)
			if accessToken == "" {
				return mr.ErrHttpClientUnauthorized.New()
			}

			currentUser, err := userProvider.MemberByToken(ctx, accessToken)
			if err != nil {
				if mr.ErrUseCaseAccessForbidden.Is(err) {
					return mr.ErrHttpAccessForbidden.New()
				}

				return err
			}

			logger.Debug(ctx, "current user", "userId", currentUser.ID().String())

			userRights := userGroups.Rights(currentUser.Group())

			if privilege != PrivilegePublic && !userRights.CheckPrivilege(privilege) {
				return mr.ErrHttpAccessForbidden.New()
			}

			if !userRights.CheckPermission(permission) {
				return mr.ErrHttpAccessForbidden.New()
			}

			// замена языка переданного клиентом в заголовке Accept-Language
			// на язык, который был установлен пользователем
			if code := currentUser.LangCode(); code != "" {
				r.Header.Set(mrreq.HeaderKeyAcceptLanguage, code)
			}

			r.Header.Set(mrreq.HeaderKeyUserIDSlashGroup, currentUser.ID().String()+"/"+currentUser.Group()) // userId/realm/kind

			if err = next(w, r); err != nil {
				if mr.ErrUseCaseAccessForbidden.Is(err) {
					return mr.ErrHttpAccessForbidden.New()
				}

				// если ошибка обработчика не связана с доступом к ресурсу
				return err
			}

			return nil
		}
	}
}

// MiddlewareHandlerIdempotency - промежуточный обработчик для организации идемпотентных запросов.
func MiddlewareHandlerIdempotency(logger mrlog.Logger, provider mridempotency.Provider, sender ResponseSender) func(next HttpHandlerFunc) HttpHandlerFunc {
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
				logger.Error(r.Context(), "MiddlewareHandlerIdempotency->Store", "error", err)
			}

			return nil
		}
	}
}
