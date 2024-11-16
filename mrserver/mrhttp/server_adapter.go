package mrhttp

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrlog"
)

// Диаграмма действия таймаутов http сервера.
//
// [Idle] [Wait] [TLS handshake] [Req.Headers] [Request body] [Response] [Idle]
//        |---------https-------||--------------------------|            - ReadTimeout
//        |---------------https---------------||-----------------------| - WriteTimeout
//                                             |-----------------------| - ServeHTTP()

const (
	defaultCaption         = "MainHttpServer"
	defaultReadTimeout     = 3 * time.Second
	defaultWriteTimeout    = 5 * time.Second
	defaultShutdownTimeout = 10 * time.Second
	serverReadyTimeout     = 60 * time.Second
)

type (
	// Adapter - Адаптер http сервера.
	Adapter struct {
		caption         string
		srv             *http.Server
		shutdownTimeout time.Duration
	}
)

// NewAdapter - создаёт объект ServerAdapter.
func NewAdapter(ctx context.Context, handler http.Handler, opts ...Option) *Adapter {
	httpServer := &http.Server{
		Handler: handler,
		// MaxHeaderBytes: 16 * 1024,
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
		BaseContext: func(_ net.Listener) context.Context {
			return ctx
		},
	}

	a := &Adapter{
		caption:         defaultCaption,
		srv:             httpServer,
		shutdownTimeout: defaultShutdownTimeout,
	}

	for _, opt := range opts {
		opt(a)
	}

	return a
}

// Caption - возвращает название http сервера.
func (a *Adapter) Caption() string {
	return a.caption
}

// ReadyTimeout - возвращает максимальное время, за которое должен быть запущен сервис.
func (a *Adapter) ReadyTimeout() time.Duration {
	return serverReadyTimeout
}

// Start - запуск http сервера.
func (a *Adapter) Start(ctx context.Context, ready func()) error {
	logger := mrlog.Ctx(ctx).With().Str("server", a.caption).Logger()
	logger.Info().Msgf("Starting the server with address: %s", a.srv.Addr)

	if ready != nil {
		ready()
	}

	if err := a.srv.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			return mrcore.ErrInternal.Wrap(err)
		}
	}

	logger.Info().Msg("The server has been stopped")

	return nil
}

// Shutdown - корректная остановка http сервера.
func (a *Adapter) Shutdown(ctx context.Context) error {
	logger := mrlog.Ctx(ctx).With().Str("server", a.caption).Logger()
	logger.Info().Msg("Shutting down the server...")

	ctx, cancel := context.WithTimeout(ctx, a.shutdownTimeout)
	defer cancel()

	if err := a.srv.Shutdown(ctx); err != nil {
		return mrcore.ErrInternal.Wrap(err)
	}

	return nil
}
