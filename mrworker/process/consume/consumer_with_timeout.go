package consume

import (
	"context"
	"errors"
	"time"

	"github.com/mondegor/go-sysmess/mrerr/mr"

	"github.com/mondegor/go-webcore/mrworker"
)

type (
	// ConsumerWithTimeout - получатель сообщений с возможностью подтверждения их получения.
	ConsumerWithTimeout struct {
		base         mrworker.MessageConsumer
		readTimeout  time.Duration
		writeTimeout time.Duration
	}
)

// NewConsumerWithTimeout - создаёт объект ConsumerWithTimeout.
func NewConsumerWithTimeout(base mrworker.MessageConsumer, readTimeout, writeTimeout time.Duration) *ConsumerWithTimeout {
	return &ConsumerWithTimeout{
		base:         base,
		readTimeout:  readTimeout,
		writeTimeout: writeTimeout,
	}
}

// ReadMessages - comment method.
func (t *ConsumerWithTimeout) ReadMessages(ctx context.Context, limit uint64) (messages []any, err error) {
	ctx, cancel := context.WithTimeout(ctx, t.readTimeout)
	defer cancel()

	messages, err = t.base.ReadMessages(ctx, limit)
	if err != nil {
		return nil, t.wrapError(err)
	}

	return messages, nil
}

// CancelMessages - comment method.
func (t *ConsumerWithTimeout) CancelMessages(ctx context.Context, messages []any) error {
	ctx, cancel := context.WithTimeout(ctx, t.writeTimeout)
	defer cancel()

	return t.wrapError(t.base.CancelMessages(ctx, messages))
}

// CommitMessage - comment method.
func (t *ConsumerWithTimeout) CommitMessage(ctx context.Context, message any, preCommit func(ctx context.Context) error) error {
	ctx, cancel := context.WithTimeout(ctx, t.writeTimeout)
	defer cancel()

	return t.wrapError(t.base.CommitMessage(ctx, message, preCommit))
}

// RejectMessage - comment method.
func (t *ConsumerWithTimeout) RejectMessage(ctx context.Context, message any, causeErr error) error {
	ctx, cancel := context.WithTimeout(ctx, t.writeTimeout)
	defer cancel()

	return t.wrapError(t.base.RejectMessage(ctx, message, causeErr))
}

func (t *ConsumerWithTimeout) wrapError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, context.DeadlineExceeded) {
		return mr.ErrInternalTimeoutPeriodHasExpired.Wrap(err)
	}

	return err
}
