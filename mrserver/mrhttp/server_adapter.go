package mrhttp

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/mondegor/go-sysmess/mrerr/mr"
	"github.com/mondegor/go-sysmess/mrlog"
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
		logger          mrlog.Logger
	}
)

// NewAdapter - создаёт объект Adapter.
func NewAdapter(handler http.Handler, opts ...Option) *Adapter {
	httpServer := &http.Server{
		Handler: handler,
		// MaxHeaderBytes: 16 * 1024,
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
	}

	a := &Adapter{
		caption:         defaultCaption,
		srv:             httpServer,
		shutdownTimeout: defaultShutdownTimeout,
	}

	for _, opt := range opts {
		opt(a)
	}

	if a.logger != nil && a.caption != "" {
		a.logger = a.logger.WithAttrs("server_name", a.caption)
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
	a.log(ctx, "Starting the server with address: "+a.srv.Addr)

	ready()

	if err := a.srv.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			return mr.ErrInternal.Wrap(err)
		}
	}

	a.log(ctx, "The server has been stopped")

	return nil
}

// Shutdown - корректная остановка http сервера.
func (a *Adapter) Shutdown(ctx context.Context) error {
	a.log(ctx, "Shutting down the server...")

	ctx, cancel := context.WithTimeout(ctx, a.shutdownTimeout)
	defer cancel()

	if err := a.srv.Shutdown(ctx); err != nil {
		return mr.ErrInternal.Wrap(err)
	}

	return nil
}

func (a *Adapter) log(ctx context.Context, msg string) {
	if a.logger != nil {
		a.logger.Info(ctx, msg)
	}
}
