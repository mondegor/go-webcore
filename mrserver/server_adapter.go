package mrserver

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrlog"
)

type (
	// ServerAdapter - comment struct.
	ServerAdapter struct {
		caption         string
		server          *http.Server
		shutdownTimeout time.Duration
		listenOpts      ListenOptions
	}

	// ServerOptions - опции для создания Server.
	ServerOptions struct {
		Caption         string
		Handler         http.Handler
		ReadTimeout     time.Duration
		WriteTimeout    time.Duration
		ShutdownTimeout time.Duration
		Listen          ListenOptions
	}

	// ListenOptions - опции для создания Listen.
	ListenOptions struct {
		BindIP string
		Port   string
	}
)

// NewServerAdapter - создаёт объект ServerAdapter.
func NewServerAdapter(ctx context.Context, opts ServerOptions) *ServerAdapter {
	httpServer := &http.Server{
		Handler: opts.Handler,
		// IdleTimeout: 120 * time.Second,
		// MaxHeaderBytes: 16 * 1024,
		// ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:  opts.ReadTimeout,
		WriteTimeout: opts.WriteTimeout,
		BaseContext: func(_ net.Listener) context.Context {
			return ctx
		},
	}

	if opts.Caption == "" {
		opts.Caption = "MainHttpServer"
	}

	return &ServerAdapter{
		caption:         opts.Caption,
		server:          httpServer,
		shutdownTimeout: opts.ShutdownTimeout,
		listenOpts:      opts.Listen,
	}
}

// PrepareToStart - comment method.
func (s *ServerAdapter) PrepareToStart(ctx context.Context) (execute func() error, interrupt func(error)) {
	return func() error {
			return s.Start(ctx)
		}, func(_ error) {
			if err := s.Shutdown(ctx); err != nil {
				mrlog.Ctx(ctx).Error().Err(err).Send()
			}
		}
}

// Start - exec server Serve().
func (s *ServerAdapter) Start(ctx context.Context) error {
	logger := mrlog.Ctx(ctx).With().Str("server", s.caption).Logger()
	logger.Info().Msg("Starting the server...")

	listener, err := s.createListener(logger)
	if err != nil {
		return mrcore.ErrInternal.Wrap(fmt.Errorf("failed start listening: %w", err))
	}

	if err = s.server.Serve(listener); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			return mrcore.ErrInternal.Wrap(err)
		}
	}

	logger.Info().Msg("Stop the server listening")

	return nil
}

// Shutdown - comment method.
func (s *ServerAdapter) Shutdown(ctx context.Context) error {
	logger := mrlog.Ctx(ctx).With().Str("server", s.caption).Logger()
	logger.Info().Msg("Shutting down the server...")

	ctx, cancel := context.WithTimeout(ctx, s.shutdownTimeout)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		return mrcore.ErrInternal.Wrap(err)
	}

	logger.Info().Msg("The server has been shutdown")

	return nil
}

func (s *ServerAdapter) createListener(logger mrlog.Logger) (net.Listener, error) {
	addr := fmt.Sprintf("%s:%s", s.listenOpts.BindIP, s.listenOpts.Port)
	logger.Debug().Msgf("Listen to TCP: %s", addr)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	logger.Info().Msgf("The server is listening to port %s", addr)

	return listener, nil
}
