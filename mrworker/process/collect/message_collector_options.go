package collect

import (
	"time"
)

type (
	// Option - настройка объекта MessageCollector.
	Option[T any] func(o *options[T])

	options[T any] struct {
		collector     *MessageCollector[T]
		captionPrefix string
	}
)

// WithCaption - устанавливает опцию caption для MessageCollector.
func WithCaption[T any](value string) Option[T] {
	return func(o *options[T]) {
		o.collector.caption = value
	}
}

// WithCaptionPrefix - устанавливает опцию caption для JobWrapper.
func WithCaptionPrefix[T any](value string) Option[T] {
	return func(o *options[T]) {
		o.captionPrefix = value
	}
}

// WithReadyTimeout - устанавливает опцию readyTimeout для MessageCollector.
func WithReadyTimeout[T any](value time.Duration) Option[T] {
	return func(o *options[T]) {
		o.collector.readyTimeout = value
	}
}

// WithFlushPeriod - устанавливает опцию периода принудительной обработки накопленной порции данных.
func WithFlushPeriod[T any](value time.Duration) Option[T] {
	return func(o *options[T]) {
		o.collector.flushPeriod = value
	}
}

// WithHandlerTimeout - устанавливает опцию handlerTimeout выполнения обработчика пачек сообщений.
func WithHandlerTimeout[T any](value time.Duration) Option[T] {
	return func(o *options[T]) {
		o.collector.handlerTimeout = value
	}
}

// WithBatchSize - устанавливает опцию размера пачки сообщений, которая будет разом обработана.
func WithBatchSize[T any](value int) Option[T] {
	return func(o *options[T]) {
		o.collector.batchSize = value
	}
}

// WithWorkersCount - устанавливает опцию количества воркеров обрабатывающих пачки сообщений.
func WithWorkersCount[T any](value int) Option[T] {
	return func(o *options[T]) {
		o.collector.workersCount = value
	}
}
