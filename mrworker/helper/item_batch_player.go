package helper

import (
	"context"
	"time"

	"github.com/mondegor/go-sysmess/errors"
	"github.com/mondegor/go-sysmess/mrevent"
	"github.com/mondegor/go-sysmess/util/conv"
)

const (
	// defaultTotalLimit - лимит обработки элементов по умолчанию.
	defaultTotalLimit = 100000

	// maxTotalLimit - максимальный допустимый лимит элементов.
	maxTotalLimit = 1000000000

	// defaultDurationLimit - лимит длительности обработки по умолчанию.
	defaultDurationLimit = time.Minute

	// maxDurationLimit - максимальная допустимая длительность обработки (1 год).
	maxDurationLimit = 365 * 24 * time.Hour
)

type (
	// ItemBatchPlayer - сервис пакетной обработки элементов в очереди.
	// Выполняет циклическую обработку элементов до их исчерпания или достижения лимитов.
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

// NewItemBatchPlayer - создаёт ItemBatchPlayer с лимитами по умолчанию.
func NewItemBatchPlayer(
	handler handler,
	eventEmitter mrevent.Emitter,
) *ItemBatchPlayer {
	return newItemBatchPlayer(handler, eventEmitter, 0, 0)
}

// NewItemBatchPlayerWithTotalLimit - создаёт ItemBatchPlayer с лимитом по количеству элементов.
func NewItemBatchPlayerWithTotalLimit(
	handler handler,
	eventEmitter mrevent.Emitter,
	totalLimit int,
) *ItemBatchPlayer {
	return newItemBatchPlayer(handler, eventEmitter, totalLimit, maxDurationLimit)
}

// NewItemBatchPlayerWithDurationLimit - создаёт ItemBatchPlayer с лимитом по времени.
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

// Execute - запускает циклическую пакетную обработку элементов.
// Параметр batchSize - размер пакета для одной итерации обработки.
//
// Процесс завершается когда:
//  1. Обработчик вернул 0 элементов (нечего обрабатывать);
//  2. Обработчик вернул меньше batchSize (последняя пачка);
//  3. Достигнут totalLimit (лимит по количеству);
//  4. Истёк durationLimit (лимит по времени);
//  5. Отменён контекст;
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
			"total":        total,
			"duration_sec": time.Since(start).Seconds(),
			"batch_size":   batchSize,
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
