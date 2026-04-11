package task

import (
	"time"

	"github.com/mondegor/go-webcore/mrworker"
)

type (
	// Option - настройка объекта JobWrapper.
	Option func(o *options)

	options struct {
		job           *JobWrapper
		captionPrefix string
	}
)

// WithCaption - устанавливает название задачи в свободной форме.
// Переопределяет значение по умолчанию ("Task").
func WithCaption(value string) Option {
	return func(o *options) {
		o.job.caption = value
	}
}

// WithCaptionPrefix - устанавливает префикс для названия задачи.
// Префикс будет добавлен перед текущим названием задачи.
func WithCaptionPrefix(value string) Option {
	return func(o *options) {
		o.captionPrefix = value
	}
}

// WithStartup - устанавливает флаг запуска задачи при старте планировщика.
func WithStartup(value bool) Option {
	return func(o *options) {
		o.job.startup = value
	}
}

// WithPeriod - устанавливает периодичность запуска задачи.
func WithPeriod(value time.Duration) Option {
	return func(o *options) {
		o.job.periodStrategy = mrworker.NewStaticPeriod(value)
	}
}

// WithPeriodStrategy - устанавливает периодичность запуска задачи.
func WithPeriodStrategy(value mrworker.PeriodStrategy) Option {
	return func(o *options) {
		o.job.periodStrategy = value
	}
}

// WithTimeout - устанавливает максимальное время выполнения задачи.
func WithTimeout(value time.Duration) Option {
	return func(o *options) {
		o.job.timeout = value
	}
}

// WithSignalDo - устанавливает канал для немедленного запуска задачи.
func WithSignalDo(value <-chan struct{}) Option {
	return func(o *options) {
		o.job.signalDo = value
	}
}
