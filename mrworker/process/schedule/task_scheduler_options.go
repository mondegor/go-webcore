package schedule

import (
	"time"

	"github.com/mondegor/go-webcore/mrworker"
)

type (
	// Option - настройка объекта TaskScheduler.
	Option func(o *options)
)

// WithCaption - устанавливает опцию caption для TaskScheduler.
func WithCaption(value string) Option {
	return func(o *options) {
		if value != "" {
			o.caption = value
		}
	}
}

// WithCaptionPrefix - устанавливает опцию caption для TaskScheduler.
func WithCaptionPrefix(value string) Option {
	return func(o *options) {
		if value != "" {
			o.captionPrefix = value
		}
	}
}

// WithReadyTimeout - устанавливает опцию readyTimeout для TaskScheduler.
func WithReadyTimeout(value time.Duration) Option {
	return func(o *options) {
		if value > 0 {
			o.readyTimeout = value
		}
	}
}

// WithTasks - устанавливает опцию tasks для TaskScheduler.
func WithTasks(values ...mrworker.Task) Option {
	return func(o *options) {
		if len(values) > 0 {
			o.tasks = values
		}
	}
}
