package mrworker

import (
	"context"
	"time"
)

type (
	// MessageConsumer - предоставляет методы для чтения, подтверждения и отклонения сообщений из очереди.
	// Используется в PULL-модели обработки сообщений.
	MessageConsumer[T any] interface {
		// ReadMessages - читает пакет сообщений из очереди.
		ReadMessages(ctx context.Context, limit int) (messages []T, err error)

		// CancelMessages - отменяет обработку не обработанных сообщений.
		// Вызывается при аварийной остановке воркеров для возврата сообщений в очередь.
		CancelMessages(ctx context.Context, messages []T) error

		// CommitMessage - подтверждает успешную обработку сообщения.
		// Параметр preCommit - функция, выполняемая перед финальным подтверждением.
		CommitMessage(ctx context.Context, message T, preCommit func(ctx context.Context) error) error

		// RejectMessage - отклоняет сообщение с указанием причины ошибки.
		RejectMessage(ctx context.Context, message T, causeErr error) error
	}

	// MessageHandler - обрабатывает отдельное сообщение с предоставлением функции фиксации результата.
	// Возвращает функцию commit для подтверждения обработки.
	// Функция commit используется для атомарного подтверждения вместе с консьюмером.
	MessageHandler[T any] interface {
		Execute(ctx context.Context, message T) (commit func(ctx context.Context) error, err error)
	}

	// MessageBatchHandler - обрабатывает пакет сообщений за один вызов.
	// Используется в PUSH-модели для пакетной обработки сообщений.
	MessageBatchHandler[T any] interface {
		Execute(ctx context.Context, messages []T) error
	}

	// Task - определяет задачу для выполнения планировщиком (TaskScheduler)
	// с настройками периода, таймаута и сигнала.
	//
	// Задачи выполняются периодически или по сигналу SignalDo.
	// Если Startup() возвращает true, задача выполняется немедленно при старте планировщика.
	Task interface {
		Caption() string
		Startup() bool
		Period() time.Duration
		Timeout() time.Duration
		SignalDo() <-chan struct{}
		Job
	}

	// Job - определяет единицу работы, которая должна быть выполнена.
	Job interface {
		Do(ctx context.Context) error
	}

	// JobFunc - обёртка для представления анонимной функции в виде Job интерфейса.
	// Позволяет передавать функции как реализации Job без создания отдельного типа.
	JobFunc func(ctx context.Context) error
)

// Do - выполняет анонимную функцию, реализуя интерфейс Job.
// Позволяет использовать JobFunc как реализацию Job интерфейса.
func (f JobFunc) Do(ctx context.Context) error {
	return f(ctx)
}
