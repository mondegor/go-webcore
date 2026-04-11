package consume

import (
	"context"
	"time"

	"github.com/mondegor/go-sysmess/errors"

	"github.com/mondegor/go-webcore/mrworker"
)

type (
	// ConsumerWithTimeout - обёртка над MessageConsumer с таймаутами на операции.
	// Добавляет таймауты на чтение (readTimeout) и запись (writeTimeout) для всех операций.
	//
	// Используется для предотвращения бесконечного ожидания при операциях с очередью.
	ConsumerWithTimeout[T any] struct {
		base         mrworker.MessageConsumer[T]
		readTimeout  time.Duration
		writeTimeout time.Duration
	}
)

// NewConsumerWithTimeout - создаёт обёртку с таймаутами для консьюмера.
func NewConsumerWithTimeout[T any](
	base mrworker.MessageConsumer[T],
	readTimeout, writeTimeout time.Duration,
) *ConsumerWithTimeout[T] {
	return &ConsumerWithTimeout[T]{
		base:         base,
		readTimeout:  readTimeout,
		writeTimeout: writeTimeout,
	}
}

// ReadMessages - читает сообщения с таймаутом на операцию чтения.
func (t *ConsumerWithTimeout[T]) ReadMessages(ctx context.Context, limit int) (messages []T, err error) {
	ctx, cancel := context.WithTimeout(ctx, t.readTimeout)
	defer cancel()

	messages, err = t.base.ReadMessages(ctx, limit)
	if err != nil {
		return nil, t.wrapError(err)
	}

	return messages, nil
}

// CancelMessages - отменяет не обработанные сообщения с таймаутом на запись.
func (t *ConsumerWithTimeout[T]) CancelMessages(ctx context.Context, messages []T) error {
	ctx, cancel := context.WithTimeout(ctx, t.writeTimeout)
	defer cancel()

	return t.wrapError(
		t.base.CancelMessages(ctx, messages),
	)
}

// CommitMessage - подтверждает обработку сообщения с таймаутом на запись.
func (t *ConsumerWithTimeout[T]) CommitMessage(ctx context.Context, message T, preCommit func(ctx context.Context) error) error {
	ctx, cancel := context.WithTimeout(ctx, t.writeTimeout)
	defer cancel()

	return t.wrapError(
		t.base.CommitMessage(ctx, message, preCommit),
	)
}

// RejectMessage - отклоняет сообщение с указанием причины и таймаутом на запись.
func (t *ConsumerWithTimeout[T]) RejectMessage(ctx context.Context, message T, causeErr error) error {
	ctx, cancel := context.WithTimeout(ctx, t.writeTimeout)
	defer cancel()

	return t.wrapError(
		t.base.RejectMessage(ctx, message, causeErr),
	)
}

func (t *ConsumerWithTimeout[T]) wrapError(err error) error {
	if errors.Is(err, context.DeadlineExceeded) {
		return errors.ErrSystemTimeoutPeriodHasExpired.Wrap(err)
	}

	return err
}
