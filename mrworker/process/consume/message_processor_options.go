package consume

import (
	"time"
)

type (
	// Option - настройка объекта MessageProcessor.
	Option[T any] func(o *options[T])

	options[T any] struct {
		processor            *MessageProcessor[T]
		captionPrefix        string
		consumerReadTimeout  time.Duration
		consumerWriteTimeout time.Duration
	}
)

// WithCaption - устанавливает опцию caption для MessageProcessor.
func WithCaption[T any](value string) Option[T] {
	return func(o *options[T]) {
		o.processor.caption = value
	}
}

// WithCaptionPrefix - устанавливает префикс для названия MessageProcessor.
func WithCaptionPrefix[T any](value string) Option[T] {
	return func(o *options[T]) {
		o.captionPrefix = value
	}
}

// WithReadyTimeout - устанавливает опцию readyTimeout для MessageProcessor.
func WithReadyTimeout[T any](value time.Duration) Option[T] {
	return func(o *options[T]) {
		o.processor.readyTimeout = value
	}
}

// WithReadPeriod - устанавливает опцию периода чтения данных консьюмером, когда он в состоянии простоя.
func WithReadPeriod[T any](value time.Duration) Option[T] {
	return func(o *options[T]) {
		o.processor.readPeriod = value
	}
}

// WithConsumerTimeout - устанавливает опцию таймаута на время отмены чтения данных
// консьюмером при неожиданном завершении работы воркеров.
func WithConsumerTimeout[T any](read, write time.Duration) Option[T] {
	return func(o *options[T]) {
		o.consumerReadTimeout = read
		o.consumerWriteTimeout = write
	}
}

// WithHandlerTimeout - устанавливает опцию handlerTimeout выполнения обработчика сообщения.
func WithHandlerTimeout[T any](value time.Duration) Option[T] {
	return func(o *options[T]) {
		o.processor.handlerTimeout = value
	}
}

// WithQueueSize - устанавливает опцию размера очереди обработки сообщений.
func WithQueueSize[T any](value int) Option[T] {
	return func(o *options[T]) {
		o.processor.queueSize = value
	}
}

// WithWorkersCount - устанавливает опцию количества воркеров обрабатывающих сообщения.
func WithWorkersCount[T any](value int) Option[T] {
	return func(o *options[T]) {
		o.processor.workersCount = value
	}
}

// WithSignalExecuteHandler - устанавливает канал сигнала для исполнения обработчика.
func WithSignalExecuteHandler[T any](ch <-chan struct{}) Option[T] {
	return func(o *options[T]) {
		o.processor.signalExecute = ch
	}
}
