package collect

import (
	"time"

	"github.com/mondegor/go-webcore/mrworker"
)

type (
	// Option - настройка объекта MessageCollector.
	Option[T any] func(o *options[T])

	options[T any] struct {
		collector     *MessageCollector[T]
		captionPrefix string
	}
)

// WithCaption - устанавливает название сервиса в свободной форме.
// Переопределяет значение по умолчанию ("MessageCollector").
func WithCaption[T any](value string) Option[T] {
	return func(o *options[T]) {
		o.collector.caption = value
	}
}

// WithCaptionPrefix - устанавливает префикс для названия сервиса.
// Префикс будет добавлен перед текущим названием.
func WithCaptionPrefix[T any](value string) Option[T] {
	return func(o *options[T]) {
		o.captionPrefix = value
	}
}

// WithReadyTimeout - устанавливает максимальное время запуска сервиса.
func WithReadyTimeout[T any](value time.Duration) Option[T] {
	return func(o *options[T]) {
		o.collector.readyTimeout = value
	}
}

// WithFlushPeriod - устанавливает период принудительной отправки накопленных сообщений.
func WithFlushPeriod[T any](value time.Duration) Option[T] {
	return func(o *options[T]) {
		o.collector.flushPeriodStrategy = mrworker.NewStaticPeriod(value)
	}
}

// WithFlushPeriodStrategy - устанавливает период принудительной отправки накопленных сообщений.
func WithFlushPeriodStrategy[T any](value mrworker.PeriodStrategy) Option[T] {
	return func(o *options[T]) {
		o.collector.flushPeriodStrategy = value
	}
}

// WithHandlerTimeout - устанавливает таймаут выполнения обработчика пакета.
func WithHandlerTimeout[T any](value time.Duration) Option[T] {
	return func(o *options[T]) {
		o.collector.handlerTimeout = value
	}
}

// WithBatchSize - устанавливает размер пакета сообщений для одной обработки.
func WithBatchSize[T any](value int) Option[T] {
	return func(o *options[T]) {
		o.collector.batchSize = value
	}
}

// WithWorkersCount - устанавливает количество параллельных воркеров-обработчиков.
// Каждый воркер обрабатывает пакет сообщений в отдельной горутине.
func WithWorkersCount[T any](value int) Option[T] {
	return func(o *options[T]) {
		o.collector.workersCount = value
	}
}
