package consume

import (
	"context"
	"time"

	"github.com/mondegor/go-sysmess/errors"

	"github.com/mondegor/go-webcore/mrworker"
)

type (
	// ConsumerWithTimeout - получатель сообщений с возможностью подтверждения их получения.
	ConsumerWithTimeout[T any] struct {
		base         mrworker.MessageConsumer[T]
		readTimeout  time.Duration
		writeTimeout time.Duration
	}
)

// NewConsumerWithTimeout - создаёт объект ConsumerWithTimeout.
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

// ReadMessages - comment method.
func (t *ConsumerWithTimeout[T]) ReadMessages(ctx context.Context, limit int) (messages []T, err error) {
	ctx, cancel := context.WithTimeout(ctx, t.readTimeout)
	defer cancel()

	messages, err = t.base.ReadMessages(ctx, limit)
	if err != nil {
		return nil, t.wrapError(err)
	}

	return messages, nil
}

// CancelMessages - comment method.
func (t *ConsumerWithTimeout[T]) CancelMessages(ctx context.Context, messages []T) error {
	ctx, cancel := context.WithTimeout(ctx, t.writeTimeout)
	defer cancel()

	return t.wrapError(
		t.base.CancelMessages(ctx, messages),
	)
}

// CommitMessage - comment method.
func (t *ConsumerWithTimeout[T]) CommitMessage(ctx context.Context, message T, preCommit func(ctx context.Context) error) error {
	ctx, cancel := context.WithTimeout(ctx, t.writeTimeout)
	defer cancel()

	return t.wrapError(
		t.base.CommitMessage(ctx, message, preCommit),
	)
}

// RejectMessage - comment method.
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
