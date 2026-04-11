package onstartup

import "time"

type (
	// Option - настройка объекта Process.
	Option func(o *options)

	options struct {
		process *Process
	}
)

// WithCaption - устанавливает название сервиса в свободной форме.
// Переопределяет значение по умолчанию ("OnStartup").
func WithCaption(value string) Option {
	return func(o *options) {
		o.process.caption = value
	}
}

// WithReadyTimeout - устанавливает максимальное время запуска сервиса.
func WithReadyTimeout(value time.Duration) Option {
	return func(o *options) {
		o.process.readyTimeout = value
	}
}
