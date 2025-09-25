package task

import "time"

type (
	// Option - настройка объекта JobWrapper.
	Option func(o *options)
)

// WithCaption - устанавливает опцию caption для JobWrapper.
func WithCaption(value string) Option {
	return func(o *options) {
		if value != "" {
			o.caption = value
		}
	}
}

// WithCaptionPrefix - устанавливает опцию caption для JobWrapper.
func WithCaptionPrefix(value string) Option {
	return func(o *options) {
		if value != "" {
			o.captionPrefix = value
		}
	}
}

// WithStartup - устанавливает опцию startup для JobWrapper.
func WithStartup(value bool) Option {
	return func(o *options) {
		o.startup = value
	}
}

// WithPeriod - устанавливает опцию period для JobWrapper.
func WithPeriod(value time.Duration) Option {
	return func(o *options) {
		if value > 0 {
			o.period = value
		}
	}
}

// WithTimeout - устанавливает опцию timeout для JobWrapper.
func WithTimeout(value time.Duration) Option {
	return func(o *options) {
		if value > 0 {
			o.timeout = value
		}
	}
}

// WithSignalDo - устанавливает опцию signalDo для JobWrapper.
func WithSignalDo(value <-chan struct{}) Option {
	return func(o *options) {
		o.signalDo = value
	}
}
