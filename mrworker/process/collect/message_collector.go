package collect

import (
	"context"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"

	"github.com/mondegor/go-sysmess/errors"
	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-sysmess/mrtrace"

	"github.com/mondegor/go-webcore/mrworker"
)

const (
	defaultCaption        = "MessageCollector"
	defaultReadyTimeout   = 30 * time.Second
	defaultFlushPeriod    = 60 * time.Second
	defaultHandlerTimeout = 30 * time.Second
	defaultBatchSize      = 100
	defaultWorkersCount   = 1
)

type (
	// MessageCollector - многопоточный сервис сбора сообщений на основе
	// внешней очереди и обработчика (PUSH модель).
	MessageCollector[T any] struct {
		caption        string
		readyTimeout   time.Duration
		flushPeriod    time.Duration
		handlerTimeout time.Duration
		batchSize      int
		workersCount   int

		handler      mrworker.MessageBatchHandler[T]
		errorHandler errors.Handler
		logger       mrlog.Logger
		traceManager mrtrace.ContextManager

		wg            sync.WaitGroup
		isSendStopped atomic.Bool
		messageQueue  chan T
		workersQueue  chan func(ctx context.Context)
		done          chan struct{}
	}
)

var (
	errInternalWorkersAreStopped       = errors.NewInternalProto("the message collector workers has been stopped")
	errInternalMessageReceptionStopped = errors.NewInternalProto("message reception in the message collector has been stopped")
)

// NewMessageCollector - создаёт объект MessageCollector.
func NewMessageCollector[T any](
	handler mrworker.MessageBatchHandler[T],
	errorHandler errors.Handler,
	logger mrlog.Logger,
	traceManager mrtrace.ContextManager,
	opts ...Option[T],
) *MessageCollector[T] {
	o := options[T]{
		collector: &MessageCollector[T]{
			caption:        defaultCaption,
			readyTimeout:   defaultReadyTimeout,
			flushPeriod:    defaultFlushPeriod,
			handlerTimeout: defaultHandlerTimeout,

			handler:      handler,
			errorHandler: errorHandler,
			logger:       logger,
			traceManager: traceManager,

			wg:           sync.WaitGroup{},
			messageQueue: make(chan T),
			workersQueue: make(chan func(ctx context.Context)),
			done:         make(chan struct{}),
		},
	}

	for _, opt := range opts {
		opt(&o)
	}

	if o.captionPrefix != "" {
		o.collector.caption = o.captionPrefix + o.collector.caption
	}

	if o.collector.batchSize < 1 {
		o.collector.batchSize = defaultBatchSize
	}

	if o.collector.workersCount < 1 {
		o.collector.workersCount = defaultWorkersCount
	}

	return o.collector
}

// Caption - возвращает название сервиса обработки сообщений.
func (p *MessageCollector[T]) Caption() string {
	return p.caption
}

// ReadyTimeout - возвращает максимальное время, за которое должен быть запущен сервис.
func (p *MessageCollector[T]) ReadyTimeout() time.Duration {
	return p.readyTimeout
}

// Start - запуск сервиса обработки сообщений.
// Отмена внешнего контекста приведёт к аварийному завершению процесса,
// для корректной остановки следует использовать Shutdown.
// Повторный запуск метода одно и того же объекта не предусмотрен, даже после вызова Shutdown.
func (p *MessageCollector[T]) Start(ctx context.Context, ready func()) error {
	p.wg.Add(1)
	defer p.wg.Done()

	p.logger.Debug(ctx, "Starting the message collector...", "collector_name", p.caption)
	defer p.logger.Debug(ctx, "The message collector has been stopped")

	wgWorkers := sync.WaitGroup{}
	workersStopped := make(chan struct{})
	ticker := time.NewTicker(p.flushPeriod)

	p.startWorkers(ctx, &wgWorkers)

	go func() {
		wgWorkers.Wait()
		close(workersStopped)
	}()

	defer func() {
		ticker.Stop()
		close(p.workersQueue)
		<-workersStopped
		close(p.messageQueue) // ??????????????
	}()

	messageBatch := make([]T, 0, p.batchSize)

	if ready != nil {
		ready()
	}

	for {
		select {
		case <-p.done:
			for {
				// в этом месте приёма новых данных уже нет,
				// но в очереди ещё могут оставаться данные, которые нужно обработать
				select {
				case message := <-p.messageQueue:
					messageBatch = append(messageBatch, message)

					if len(messageBatch) < p.batchSize {
						continue
					}
				default:
				}

				break
			}
		case <-ctx.Done():
			p.logger.Debug(ctx, "The message collector detected context 'Done'", "error", ctx.Err())

			// предварительно завершается приём данных
			p.isSendStopped.Store(true)

			// принудительная очистка очереди,
			// т.к. контекст отменён и данные уже не получится обработать
			for {
				select {
				case <-p.messageQueue:
				default:
					return nil
				}
			}
		case message := <-p.messageQueue:
			messageBatch = append(messageBatch, message)

			if len(messageBatch) < p.batchSize {
				continue
			}
		case <-ticker.C:
			p.logger.Debug(ctx, "The message collector ticker.C")
		}

		if len(messageBatch) == 0 {
			if p.isSendStopped.Load() {
				return nil // если данных нет и их приём остановлен, то процесс завершается
			}

			continue
		}

		p.logger.Info(ctx, "Got message batch in the message collector...", "message_batch", len(messageBatch))

		select {
		case <-workersStopped:
			return errInternalWorkersAreStopped.New("collector_name", p.caption)
		case p.workersQueue <- p.workerFunc(messageBatch):
			messageBatch = make([]T, 0, p.batchSize)
		}
	}
}

// PushMessage - отправляет сообщение в очередь для дальнейшей её обработки.
func (p *MessageCollector[T]) PushMessage(ctx context.Context, message T) error {
	if p.isSendStopped.Load() {
		return errInternalMessageReceptionStopped.New("collector_name", p.caption)
	}

	select {
	case p.messageQueue <- message:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Shutdown - корректная остановка сервиса обработки сообщений.
// При повторном вызове метода произойдёт panic.
func (p *MessageCollector[T]) Shutdown(ctx context.Context) error {
	p.logger.Debug(ctx, "Shutting down the message collector...")
	p.isSendStopped.Store(true) // завершается приём данных
	close(p.done)

	p.wg.Wait()
	p.logger.Debug(ctx, "The message collector has been shut down")

	return nil
}

func (p *MessageCollector[T]) startWorkers(ctx context.Context, wg *sync.WaitGroup) {
	for i := 0; i < p.workersCount; i++ {
		wg.Add(1)

		go func(ctx context.Context) {
			ctx = p.traceManager.WithGeneratedProcessID(ctx, mrtrace.KeyWorkerID)

			defer func() {
				wg.Done()

				if rvr := recover(); rvr != nil {
					p.errorHandler.Handle(
						ctx,
						errors.ErrInternalCaughtPanic.New(
							"source", "message collector: "+p.caption,
							"recover", rvr,
							"stack_trace", string(debug.Stack()),
						),
					)
				}
			}()

			for fn := range p.workersQueue {
				fn(ctx)
			}

			p.logger.Debug(ctx, "The worker has been stopped")
		}(ctx)
	}
}

func (p *MessageCollector[T]) workerFunc(messages []T) func(ctx context.Context) {
	return func(ctx context.Context) {
		handlerCtx, cancel := context.WithTimeout(p.traceManager.WithGeneratedProcessID(ctx, mrtrace.KeyTaskID), p.handlerTimeout)
		defer cancel()

		p.logger.Debug(ctx, "workerFunc", "message_batch", len(messages), "message[0]", messages[0])

		if err := p.handler.Execute(handlerCtx, messages); err != nil {
			p.errorHandler.Handle(ctx, err)

			return
		}

		p.logger.Debug(ctx, "The handler has been successfully executed")
	}
}
