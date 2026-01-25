package collect

import (
	"time"
)

type (
	// Option - настройка объекта MessageCollector.
	Option func(o *options)

	options struct {
		collector     *MessageCollector
		captionPrefix string
	}
)

// WithCaption - устанавливает опцию caption для MessageCollector.
func WithCaption(value string) Option {
	return func(o *options) {
		o.collector.caption = value
	}
}

// WithCaptionPrefix - устанавливает опцию caption для JobWrapper.
func WithCaptionPrefix(value string) Option {
	return func(o *options) {
		o.captionPrefix = value
	}
}

// WithReadyTimeout - устанавливает опцию readyTimeout для MessageCollector.
func WithReadyTimeout(value time.Duration) Option {
	return func(o *options) {
		o.collector.readyTimeout = value
	}
}

// WithFlushPeriod - устанавливает опцию периода принудительной обработки накопленной порции данных.
func WithFlushPeriod(value time.Duration) Option {
	return func(o *options) {
		o.collector.flushPeriod = value
	}
}

// WithHandlerTimeout - устанавливает опцию handlerTimeout выполнения обработчика пачек сообщений.
func WithHandlerTimeout(value time.Duration) Option {
	return func(o *options) {
		o.collector.handlerTimeout = value
	}
}

// WithBatchSize - устанавливает опцию размера пачки сообщений, которая будет разом обработана.
func WithBatchSize(value int) Option {
	return func(o *options) {
		o.collector.batchSize = value
	}
}

// WithWorkersCount - устанавливает опцию количества воркеров обрабатывающих пачки сообщений.
func WithWorkersCount(value int) Option {
	return func(o *options) {
		o.collector.workersCount = value
	}
}
