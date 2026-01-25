package httpserver

import (
	"context"
	"net"
	"time"

	"github.com/mondegor/go-sysmess/mrlog"
)

type (
	// Option - настройка http сервера Adapter.
	Option func(o *options)

	options struct {
		server *Adapter
	}
)

// WithHostPort - устанавливает хост (IP) и порт на которых будет работать http сервер.
func WithHostPort(host, port string) Option {
	return func(o *options) {
		if port != "" {
			port = ":" + port
		}

		o.server.srv.Addr = host + port
	}
}

// WithLogger - устанавливает логгер для логирования работы сервера.
func WithLogger(value mrlog.Logger) Option {
	return func(o *options) {
		o.server.logger = value
	}
}

// WithBaseContext - устанавливает контекст, который будет использоваться в каждом запросе.
func WithBaseContext(ctx context.Context) Option {
	return func(o *options) {
		o.server.srv.BaseContext = func(_ net.Listener) context.Context {
			return ctx
		}
	}
}

// WithCaption - устанавливает название сервера.
func WithCaption(value string) Option {
	return func(o *options) {
		o.server.caption = value
	}
}

// WithReadTimeout - устанавливает таймаут при чтении заголовка и тела запроса.
func WithReadTimeout(value time.Duration) Option {
	return func(o *options) {
		o.server.srv.ReadTimeout = value
	}
}

// WithWriteTimeout - устанавливает таймаут на формирование ответа сервера.
func WithWriteTimeout(value time.Duration) Option {
	return func(o *options) {
		o.server.srv.WriteTimeout = value
	}
}

// WithShutdownTimeout - устанавливает таймаут для корректного
// завершения активных соединений при остановке сервера.
func WithShutdownTimeout(value time.Duration) Option {
	return func(o *options) {
		o.server.shutdownTimeout = value
	}
}
