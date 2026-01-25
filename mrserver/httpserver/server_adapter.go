package httpserver

import (
	"context"
	"net/http"
	"time"

	"github.com/mondegor/go-sysmess/errors"
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

// New - создаёт объект Adapter.
func New(handler http.Handler, opts ...Option) *Adapter {
	o := options{
		server: &Adapter{
			caption: defaultCaption,
			srv: &http.Server{
				Handler: handler,
				// MaxHeaderBytes: 16 * 1024,
				ReadTimeout:  defaultReadTimeout,
				WriteTimeout: defaultWriteTimeout,
			},
			shutdownTimeout: defaultShutdownTimeout,
		},
	}

	for _, opt := range opts {
		opt(&o)
	}

	if o.server.logger != nil && o.server.caption != "" {
		o.server.logger = mrlog.WithAttrs(o.server.logger, "server_name", o.server.caption)
	}

	return o.server
}

// Caption - возвращает название http сервера.
func (s *Adapter) Caption() string {
	return s.caption
}

// ReadyTimeout - возвращает максимальное время, за которое должен быть запущен сервис.
func (s *Adapter) ReadyTimeout() time.Duration {
	return serverReadyTimeout
}

// Start - запуск http сервера.
func (s *Adapter) Start(ctx context.Context, ready func()) error {
	s.log(ctx, "Starting the server with address: "+s.srv.Addr)

	ready()

	if err := s.srv.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			return errors.WrapInternalError(err, "listening server failed")
		}
	}

	s.log(ctx, "The server has been stopped")

	return nil
}

// Shutdown - корректная остановка http сервера.
func (s *Adapter) Shutdown(ctx context.Context) error {
	s.log(ctx, "Shutting down the server...")

	ctx, cancel := context.WithTimeout(ctx, s.shutdownTimeout)
	defer cancel()

	if err := s.srv.Shutdown(ctx); err != nil {
		return errors.WrapInternalError(err, "shutting down server failed")
	}

	return nil
}

func (s *Adapter) log(ctx context.Context, msg string) {
	if s.logger != nil {
		s.logger.Info(ctx, msg)
	}
}
