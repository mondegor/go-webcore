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

// WithHostPort - устанавливает адрес (хост и порт) для прослушивания HTTP-сервером.
func WithHostPort(host, port string) Option {
	return func(o *options) {
		if port != "" {
			port = ":" + port
		}

		o.server.srv.Addr = host + port
	}
}

// WithLogger - устанавливает логгер для записи событий запуска и остановки сервера.
func WithLogger(value mrlog.Logger) Option {
	return func(o *options) {
		o.server.logger = value
	}
}

// WithBaseContext - устанавливает базовый контекст для всех HTTP-запросов.
// Контекст будет использоваться как родительский для каждого входящего запроса.
func WithBaseContext(ctx context.Context) Option {
	return func(o *options) {
		o.server.srv.BaseContext = func(_ net.Listener) context.Context {
			return ctx
		}
	}
}

// WithCaption - устанавливает название сервера в свободной форме.
// Переопределяет значение по умолчанию ("MainHttpServer").
func WithCaption(value string) Option {
	return func(o *options) {
		o.server.caption = value
	}
}

// WithReadTimeout - устанавливает максимальное время чтения заголовка и тела запроса.
func WithReadTimeout(value time.Duration) Option {
	return func(o *options) {
		o.server.srv.ReadTimeout = value
	}
}

// WithWriteTimeout - устанавливает максимальное время формирования ответа.
func WithWriteTimeout(value time.Duration) Option {
	return func(o *options) {
		o.server.srv.WriteTimeout = value
	}
}

// WithShutdownTimeout - устанавливает таймаут корректной остановки сервера.
// Если таймаут истёк, активные соединения будут принудительно закрыты.
func WithShutdownTimeout(value time.Duration) Option {
	return func(o *options) {
		o.server.shutdownTimeout = value
	}
}
