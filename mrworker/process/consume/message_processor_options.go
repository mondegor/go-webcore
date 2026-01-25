package consume

import (
	"time"
)

type (
	// Option - настройка объекта MessageProcessor.
	Option func(o *options)

	options struct {
		processor            *MessageProcessor
		captionPrefix        string
		consumerReadTimeout  time.Duration
		consumerWriteTimeout time.Duration
	}
)

// WithCaption - устанавливает опцию caption для MessageProcessor.
func WithCaption(value string) Option {
	return func(o *options) {
		o.processor.caption = value
	}
}

// WithCaptionPrefix - устанавливает опцию caption для JobWrapper.
func WithCaptionPrefix(value string) Option {
	return func(o *options) {
		o.captionPrefix = value
	}
}

// WithReadyTimeout - устанавливает опцию readyTimeout для MessageProcessor.
func WithReadyTimeout(value time.Duration) Option {
	return func(o *options) {
		o.processor.readyTimeout = value
	}
}

// WithReadPeriod - устанавливает опцию периода чтения данных консьюмером, когда он в состоянии простоя.
func WithReadPeriod(value time.Duration) Option {
	return func(o *options) {
		o.processor.readPeriod = value
	}
}

// WithConsumerTimeout - устанавливает опцию таймаута на время отмены чтения данных
// консьюмером при неожиданном завершении работы воркеров.
func WithConsumerTimeout(read, write time.Duration) Option {
	return func(o *options) {
		o.consumerReadTimeout = read
		o.consumerWriteTimeout = write
	}
}

// WithHandlerTimeout - устанавливает опцию handlerTimeout выполнения обработчика сообщения.
func WithHandlerTimeout(value time.Duration) Option {
	return func(o *options) {
		o.processor.handlerTimeout = value
	}
}

// WithQueueSize - устанавливает опцию размера очереди обработки сообщений.
func WithQueueSize(value int) Option {
	return func(o *options) {
		o.processor.queueSize = value
	}
}

// WithWorkersCount - устанавливает опцию количества воркеров обрабатывающих сообщения.
func WithWorkersCount(value int) Option {
	return func(o *options) {
		o.processor.workersCount = value
	}
}

// WithSignalExecuteHandler - comment func.
func WithSignalExecuteHandler(ch <-chan struct{}) Option {
	return func(o *options) {
		o.processor.signalExecute = ch
	}
}
