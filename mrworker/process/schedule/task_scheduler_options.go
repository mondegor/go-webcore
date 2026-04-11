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

// WithCaption - устанавливает название планировщика в свободной форме.
// Переопределяет значение по умолчанию ("TaskScheduler").
func WithCaption(value string) Option {
	return func(o *options) {
		o.scheduler.caption = value
	}
}

// WithCaptionPrefix - устанавливает префикс для названия планировщика.
// Префикс будет добавлен перед текущим названием.
func WithCaptionPrefix(value string) Option {
	return func(o *options) {
		o.captionPrefix = value
	}
}

// WithReadyTimeout - устанавливает максимальное время запуска планировщика.
func WithReadyTimeout(value time.Duration) Option {
	return func(o *options) {
		o.scheduler.readyTimeout = value
	}
}

// WithTasks - размещает задачи в планировщике.
// Можно вызывать несколько раз для добавления задач группами.
func WithTasks(values ...mrworker.Task) Option {
	return func(o *options) {
		o.scheduler.tasks = append(o.scheduler.tasks, values...)
	}
}
