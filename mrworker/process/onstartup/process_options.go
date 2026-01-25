package onstartup

import "time"

type (
	// Option - настройка объекта Process.
	Option func(o *options)

	options struct {
		process *Process
	}
)

// WithCaption - устанавливает опцию caption для Process.
func WithCaption(value string) Option {
	return func(o *options) {
		o.process.caption = value
	}
}

// WithReadyTimeout - устанавливает опцию readyTimeout для Process.
func WithReadyTimeout(value time.Duration) Option {
	return func(o *options) {
		o.process.readyTimeout = value
	}
}
