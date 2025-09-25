package collect

import (
	"time"
)

type (
	// Option - настройка объекта MessageCollector.
	Option func(o *options)
)

// WithCaption - устанавливает опцию caption для MessageCollector.
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

// WithReadyTimeout - устанавливает опцию readyTimeout для MessageCollector.
func WithReadyTimeout(value time.Duration) Option {
	return func(o *options) {
		if value > 0 {
			o.readyTimeout = value
		}
	}
}

// WithFlushPeriod - устанавливает опцию периода принудительной обработки накопленной порции данных.
func WithFlushPeriod(value time.Duration) Option {
	return func(o *options) {
		if value > 0 {
			o.flushPeriod = value
		}
	}
}

// WithHandlerTimeout - устанавливает опцию handlerTimeout выполнения обработчика пачек сообщений.
func WithHandlerTimeout(value time.Duration) Option {
	return func(o *options) {
		if value > 0 {
			o.handlerTimeout = value
		}
	}
}

// WithBatchSize - устанавливает опцию размера пачки сообщений, которая будет разом обработана.
func WithBatchSize(value uint64) Option {
	return func(o *options) {
		if o.batchSize > 0 {
			o.batchSize = value
		}
	}
}

// WithWorkersCount - устанавливает опцию количества воркеров обрабатывающих пачки сообщений.
func WithWorkersCount(value uint64) Option {
	return func(o *options) {
		if value > 0 {
			o.workersCount = value
		}
	}
}
