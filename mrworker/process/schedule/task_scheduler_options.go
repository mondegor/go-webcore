package schedule

import (
	"time"

	"github.com/mondegor/go-webcore/mrworker"
)

type (
	// Option - настройка объекта TaskScheduler.
	Option func(o *options)

	options struct {
		scheduler     *TaskScheduler
		captionPrefix string
	}
)

// WithCaption - устанавливает опцию caption для TaskScheduler.
func WithCaption(value string) Option {
	return func(o *options) {
		o.scheduler.caption = value
	}
}

// WithCaptionPrefix - устанавливает префикс для названия TaskScheduler.
func WithCaptionPrefix(value string) Option {
	return func(o *options) {
		o.captionPrefix = value
	}
}

// WithReadyTimeout - устанавливает опцию readyTimeout для TaskScheduler.
func WithReadyTimeout(value time.Duration) Option {
	return func(o *options) {
		o.scheduler.readyTimeout = value
	}
}

// WithTasks - устанавливает опцию tasks для TaskScheduler.
func WithTasks(values ...mrworker.Task) Option {
	return func(o *options) {
		o.scheduler.tasks = append(o.scheduler.tasks, values...)
	}
}
