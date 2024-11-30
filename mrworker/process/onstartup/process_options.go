package onstartup

import "time"

type (
	// Option - настройка объекта Process.
	Option func(p *Process)
)

// WithCaption - устанавливает опцию caption для Process.
func WithCaption(value string) Option {
	return func(p *Process) {
		if value != "" {
			p.caption = value
		}
	}
}

// WithReadyTimeout - устанавливает опцию readyTimeout для Process.
func WithReadyTimeout(value time.Duration) Option {
	return func(p *Process) {
		if value > 0 {
			p.readyTimeout = value
		}
	}
}
