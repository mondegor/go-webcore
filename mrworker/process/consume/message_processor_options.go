package consume

import (
	"time"
)

type (
	// Option - настройка объекта MessageProcessor.
	Option func(o *options)
)

// WithCaption - устанавливает опцию caption для MessageProcessor.
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

// WithReadyTimeout - устанавливает опцию readyTimeout для MessageProcessor.
func WithReadyTimeout(value time.Duration) Option {
	return func(o *options) {
		if value > 0 {
			o.readyTimeout = value
		}
	}
}

// WithReadPeriod - устанавливает опцию периода чтения данных консьюмером, когда он в состоянии простоя.
func WithReadPeriod(value time.Duration) Option {
	return func(o *options) {
		if value > 0 {
			o.readPeriod = value
		}
	}
}

// WithConsumerTimeout - устанавливает опцию таймаута на время отмены чтения данных
// консьюмером при неожиданном завершении работы воркеров.
func WithConsumerTimeout(read, write time.Duration) Option {
	return func(o *options) {
		if read > 0 {
			o.consumerReadTimeout = read
		}

		if write > 0 {
			o.consumerWriteTimeout = write
		}
	}
}

// WithHandlerTimeout - устанавливает опцию handlerTimeout выполнения обработчика сообщения.
func WithHandlerTimeout(value time.Duration) Option {
	return func(o *options) {
		if value > 0 {
			o.handlerTimeout = value
		}
	}
}

// WithQueueSize - устанавливает опцию размера очереди обработки сообщений.
func WithQueueSize(value uint64) Option {
	return func(o *options) {
		if o.queueSize > 0 {
			o.queueSize = value
		}
	}
}

// WithWorkersCount - устанавливает опцию количества воркеров обрабатывающих сообщения.
func WithWorkersCount(value uint64) Option {
	return func(o *options) {
		if value > 0 {
			o.workersCount = value
		}
	}
}

// WithSignalExecuteHandler - comment func.
func WithSignalExecuteHandler(ch <-chan struct{}) Option {
	return func(o *options) {
		o.signalExecute = ch
	}
}
