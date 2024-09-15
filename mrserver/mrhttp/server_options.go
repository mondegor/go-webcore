package mrhttp

import (
	"time"
)

// Option - опция используемая при создании http сервера.
type (
	Option func(s *Adapter)
)

// WithHostAndPort - устанавливает хост (IP) и порт на которых будет работать http сервер.
func WithHostAndPort(host, port string) Option {
	return func(s *Adapter) {
		if port != "" {
			port = ":" + port
		}

		s.srv.Addr = host + port
	}
}

// WithCaption - устанавливает название сервера.
func WithCaption(value string) Option {
	return func(s *Adapter) {
		s.caption = value
	}
}

// WithReadTimeout - устанавливает таймаут при чтении заголовка и тела запроса.
func WithReadTimeout(value time.Duration) Option {
	return func(s *Adapter) {
		s.srv.ReadTimeout = value
	}
}

// WithWriteTimeout - устанавливает таймаут на формирование ответа сервера.
func WithWriteTimeout(value time.Duration) Option {
	return func(s *Adapter) {
		s.srv.WriteTimeout = value
	}
}

// WithShutdownTimeout - устанавливает таймаут для корректного
// завершения активных соединений при остановке сервера.
func WithShutdownTimeout(value time.Duration) Option {
	return func(s *Adapter) {
		s.shutdownTimeout = value
	}
}
