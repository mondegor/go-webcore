package helper

import (
	"context"
	"time"

	"github.com/mondegor/go-sysmess/errors"
	"github.com/mondegor/go-sysmess/mrevent"
	"github.com/mondegor/go-sysmess/util/conv"
)

const (
	defaultTotalLimit    = 100000
	maxTotalLimit        = 1000000000
	defaultDurationLimit = time.Minute
	maxDurationLimit     = 365 * 24 * time.Hour
)

type (
	// ItemBatchPlayer - объект изменяющий статусы сломавшихся элементов, находящихся в очереди.
	ItemBatchPlayer struct {
		handler       handler
		eventEmitter  mrevent.Emitter
		totalLimit    int
		durationLimit time.Duration
	}

	handler interface {
		Execute(ctx context.Context, limit int) (count int, err error)
	}
)

// NewItemBatchPlayer - создаёт объект ItemBatchPlayer.
func NewItemBatchPlayer(
	handler handler,
	eventEmitter mrevent.Emitter,
) *ItemBatchPlayer {
	return newItemBatchPlayer(handler, eventEmitter, 0, 0)
}

// NewItemBatchPlayerWithTotalLimit - создаёт объект ItemBatchPlayer.
func NewItemBatchPlayerWithTotalLimit(
	handler handler,
	eventEmitter mrevent.Emitter,
	totalLimit int,
) *ItemBatchPlayer {
	return newItemBatchPlayer(handler, eventEmitter, totalLimit, maxDurationLimit)
}

// NewItemBatchPlayerWithDurationLimit - создаёт объект ItemBatchPlayer.
func NewItemBatchPlayerWithDurationLimit(
	handler handler,
	eventEmitter mrevent.Emitter,
	durationLimit time.Duration,
) *ItemBatchPlayer {
	return newItemBatchPlayer(handler, eventEmitter, maxTotalLimit, durationLimit)
}

func newItemBatchPlayer(
	handler handler,
	eventEmitter mrevent.Emitter,
	totalLimit int,
	durationLimit time.Duration,
) *ItemBatchPlayer {
	if totalLimit <= 0 {
		totalLimit = defaultTotalLimit
	}

	if durationLimit <= 0 {
		durationLimit = defaultDurationLimit
	}

	return &ItemBatchPlayer{
		handler:       handler,
		eventEmitter:  eventEmitter,
		totalLimit:    totalLimit,
		durationLimit: durationLimit,
	}
}

// Execute - запускает в цикле вложенный процесс пакетной обработки элементов.
// Процесс завершается когда элементов для обработки не остаётся или
// превышен лимит кол-ва обработанных элементов.
func (p *ItemBatchPlayer) Execute(ctx context.Context, batchSize int) error {
	if batchSize < 1 {
		return errors.ErrInternalIncorrectInputData.WithDetails("batchSize is zero or negative")
	}

	if batchSize > p.totalLimit {
		return errors.ErrInternalIncorrectInputData.WithDetails("batchSize > totalLimit")
	}

	total := 0
	start := time.Now()

	for {
		count, err := p.handler.Execute(ctx, batchSize)
		if err != nil {
			return p.wrapError(err)
		}

		select {
		case <-ctx.Done(): // на случай, если обработчик не обработал контекст
			return p.wrapError(ctx.Err())
		default:
		}

		total += count

		if count == 0 ||
			batchSize > count ||
			total >= p.totalLimit ||
			time.Since(start) >= p.durationLimit {
			break
		}
	}

	p.eventEmitter.Emit(
		ctx,
		"Execute",
		conv.Group{
			"total":           total,
			"durationSeconds": time.Since(start).Seconds(),
			"batchSize":       batchSize,
		},
	)

	return nil
}

func (p *ItemBatchPlayer) wrapError(err error) error {
	if errors.Is(err, context.DeadlineExceeded) {
		return errors.ErrSystemTimeoutPeriodHasExpired.Wrap(err)
	}

	return err
}
