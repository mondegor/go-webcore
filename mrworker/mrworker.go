package mrworker

import (
	"context"
	"time"
)

type (
	// MessageConsumer - предоставляет методы для чтения, подтверждения и отклонения сообщений из очереди.
	MessageConsumer[T any] interface {
		ReadMessages(ctx context.Context, limit int) (messages []T, err error)
		CancelMessages(ctx context.Context, messages []T) error
		CommitMessage(ctx context.Context, message T, preCommit func(ctx context.Context) error) error
		RejectMessage(ctx context.Context, message T, causeErr error) error
	}

	// MessageHandler - обрабатывает отдельное сообщение
	// с предоставлением функции фиксации результата.
	MessageHandler[T any] interface {
		Execute(ctx context.Context, message T) (commit func(ctx context.Context) error, err error)
	}

	// MessageBatchHandler - обрабатывает пакет сообщений за один вызов.
	MessageBatchHandler[T any] interface {
		Execute(ctx context.Context, messages []T) error
	}

	// Task - определяет задачу для выполнения планировщиком
	// с настройками периода, таймаута и сигнала.
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

	// JobFunc - обёртка для возможности представления анонимной функции в виде Job интерфейса.
	JobFunc func(ctx context.Context) error
)

// Do - запускает анонимную функцию в Job интерфейсе.
func (f JobFunc) Do(ctx context.Context) error {
	return f(ctx)
}
