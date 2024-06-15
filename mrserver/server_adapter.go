package mrserver

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"path/filepath"
	"time"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrlog"
)

const (
	ListenTypeSock = "sock" // ListenTypeSock - тип прослушивания: сокет
	ListenTypePort = "port" // ListenTypePort - тип прослушивания: IP + порт
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
		AppPath  string
		Type     string
		SockName string
		BindIP   string
		Port     string
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
		return mrcore.ErrInternal.Wrap(err)
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
	if s.listenOpts.Type == ListenTypeSock {
		logger.Debug().Msg("Detect app real path")

		appDir, err := filepath.Abs(filepath.Dir(s.listenOpts.AppPath))
		if err != nil {
			return nil, fmt.Errorf("app real path: %w", err)
		}

		socketPath := filepath.Join(appDir, s.listenOpts.SockName)
		logger.Debug().Msgf("Listen to unix socket: %s", socketPath)

		listener, err := net.Listen("unix", socketPath)

		if err == nil {
			logger.Info().Msgf("The server is listening unix socket: %s", socketPath)
		}

		return listener, err
	} else if s.listenOpts.Type == ListenTypePort {
		addr := fmt.Sprintf("%s:%s", s.listenOpts.BindIP, s.listenOpts.Port)
		logger.Debug().Msgf("Listen to tcp: %s", addr)

		listener, err := net.Listen("tcp", addr)

		if err == nil {
			logger.Info().Msgf("The server is listening to port %s", addr)
		}

		return listener, err
	}

	availableValues := fmt.Sprintf("Available values: %s, %s", ListenTypePort, ListenTypeSock)

	if s.listenOpts.Type == "" {
		return nil, fmt.Errorf("listen type is required. %s", availableValues)
	}

	return nil, fmt.Errorf("listen type '%s' is unknown. %s", s.listenOpts.Type, availableValues)
}
