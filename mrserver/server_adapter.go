package mrserver

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/mondegor/go-webcore/mrcore"
)

type (
	ServerAdapter struct {
		server *http.Server
		notifyChan chan error
		shutdownTimeout time.Duration
		logger mrcore.Logger
	}

	ServerOptions struct {
		Handler http.Handler
		ReadTimeout time.Duration
		WriteTimeout time.Duration
		ShutdownTimeout time.Duration
	}
)

func NewServer(logger mrcore.Logger, opt ServerOptions) *ServerAdapter {
	httpServer := &http.Server{
		Handler: opt.Handler,
		// IdleTimeout: 120 * time.Second,
		// MaxHeaderBytes: 16 * 1024,
		// ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout: opt.ReadTimeout,
		WriteTimeout: opt.WriteTimeout,
	}

	return &ServerAdapter{
		server: httpServer,
		notifyChan: make(chan error, 1),
		shutdownTimeout: opt.ShutdownTimeout,
		logger: logger,
	}
}

func (s *ServerAdapter) Start(opt ListenOptions) error {
	listener, err := s.createListener(&opt)

	if err != nil {
		return fmt.Errorf("http server start: %w", err)
	}

	go func() {
		defer close(s.notifyChan)
		s.notifyChan <- s.server.Serve(listener)
	}()

	return nil
}

func (s *ServerAdapter) Notify() <-chan error {
	return s.notifyChan
}

func (s *ServerAdapter) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	return s.server.Shutdown(ctx)
}
