package consume

import (
	"time"

	"github.com/mondegor/go-webcore/mrworker"
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

// WithCaption - устанавливает название сервиса в свободной форме.
// Переопределяет значение по умолчанию ("MessageProcessor").
func WithCaption[T any](value string) Option[T] {
	return func(o *options[T]) {
		o.processor.caption = value
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
		o.processor.readyTimeout = value
	}
}

// WithReadPeriod - устанавливает период опроса очереди в состоянии простоя.
func WithReadPeriod[T any](value time.Duration) Option[T] {
	return func(o *options[T]) {
		o.processor.readPeriodStrategy = mrworker.NewStaticPeriodStrategy(value)
	}
}

// WithReadPeriodStrategy - устанавливает период опроса очереди в состоянии простоя.
func WithReadPeriodStrategy[T any](value mrworker.PeriodStrategy) Option[T] {
	return func(o *options[T]) {
		o.processor.readPeriodStrategy = value
	}
}

// WithConsumerTimeout - устанавливает таймауты на операции консьюмера.
// Если хотя бы одно значение > 0, консьюмер автоматически оборачивается в ConsumerWithTimeout.
func WithConsumerTimeout[T any](read, write time.Duration) Option[T] {
	return func(o *options[T]) {
		o.consumerReadTimeout = read
		o.consumerWriteTimeout = write
	}
}

// WithHandlerTimeout - устанавливает таймаут выполнения обработчика сообщения.
func WithHandlerTimeout[T any](value time.Duration) Option[T] {
	return func(o *options[T]) {
		o.processor.handlerTimeout = value
	}
}

// WithQueueSize - устанавливает максимальное количество сообщений для одной выборки.
func WithQueueSize[T any](value int) Option[T] {
	return func(o *options[T]) {
		o.processor.queueSize = value
	}
}

// WithWorkersCount - устанавливает количество параллельных воркеров-обработчиков.
// Каждый воркер обрабатывает одно сообщение в отдельной горутине.
func WithWorkersCount[T any](value int) Option[T] {
	return func(o *options[T]) {
		o.processor.workersCount = value
	}
}

// WithSignalExecuteHandler - устанавливает канал для немедленного опроса очереди.
func WithSignalExecuteHandler[T any](ch <-chan struct{}) Option[T] {
	return func(o *options[T]) {
		o.processor.signalExecute = ch
	}
}
