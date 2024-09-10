package mrhttpserver

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrlog"
)

// Диаграмма действия таймаутов сервера.
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
)

type (
	// Adapter - comment struct.
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

	adapter := &Adapter{
		caption:         defaultCaption,
		srv:             httpServer,
		shutdownTimeout: defaultShutdownTimeout,
	}

	for _, opt := range opts {
		opt(adapter)
	}

	return adapter
}

// PrepareToStart - comment method.
func (s *Adapter) PrepareToStart(ctx context.Context) (execute func() error, interrupt func(error)) {
	return func() error {
			return s.Start(ctx)
		}, func(_ error) {
			if err := s.Shutdown(ctx); err != nil {
				mrlog.Ctx(ctx).Error().Err(err).Send()
			}
		}
}

// Start - запуск http сервера.
func (s *Adapter) Start(ctx context.Context) error {
	logger := mrlog.Ctx(ctx).With().Str("server", s.caption).Logger()
	logger.Info().Msgf("Starting the server with address: %s", s.srv.Addr)

	if err := s.srv.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			return mrcore.ErrInternal.Wrap(err)
		}
	}

	logger.Info().Msg("Stop the server listening")

	return nil
}

// Shutdown - корректная остановка http сервера.
func (s *Adapter) Shutdown(ctx context.Context) error {
	logger := mrlog.Ctx(ctx).With().Str("server", s.caption).Logger()
	logger.Info().Msg("Shutting down the server...")

	ctx, cancel := context.WithTimeout(ctx, s.shutdownTimeout)
	defer cancel()

	if err := s.srv.Shutdown(ctx); err != nil {
		return mrcore.ErrInternal.Wrap(err)
	}

	logger.Info().Msg("The server has been shutdown")

	return nil
}
