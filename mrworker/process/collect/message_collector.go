package collect

import (
	"context"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"

	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-sysmess/mrerr/mr"
	"github.com/mondegor/go-sysmess/mrlog"

	"github.com/mondegor/go-webcore/mrcore"
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
	// MessageCollector - многопоточный сервис сбора сообщений на основе внешней очереди и обработчика.
	MessageCollector struct {
		caption        string
		readyTimeout   time.Duration
		flushPeriod    time.Duration
		handlerTimeout time.Duration
		batchSize      uint64
		workersCount   uint64

		handler         mrworker.MessageBatchHandler
		errorHandler    mrcore.ErrorHandler
		logger          mrlog.Logger
		contextEmbedder contextEmbedder

		wgMain        sync.WaitGroup
		isSendStopped atomic.Bool
		messageQueue  chan []byte
		workersQueue  chan func(ctx context.Context)
		done          chan struct{}
	}

	contextEmbedder interface {
		NewContextWithIDs(originalCtx context.Context) context.Context
		WithWorkerIDContext(ctx context.Context) context.Context
		WithTaskIDContext(ctx context.Context) context.Context
	}

	options struct {
		caption        string
		captionPrefix  string
		readyTimeout   time.Duration
		flushPeriod    time.Duration
		handlerTimeout time.Duration
		batchSize      uint64
		workersCount   uint64
	}
)

var (
	errInternalWorkersAreStopped       = mrerr.NewKindInternal("the message collector workers has been stopped")
	errInternalMessageReceptionStopped = mrerr.NewKindInternal("message reception in the message collector has been stopped")
)

// NewMessageCollector - создаёт объект MessageCollector.
func NewMessageCollector(
	handler mrworker.MessageBatchHandler,
	errorHandler mrcore.ErrorHandler,
	logger mrlog.Logger,
	idGenerator contextEmbedder,
	opts ...Option,
) *MessageCollector {
	o := options{
		caption:        defaultCaption,
		readyTimeout:   defaultReadyTimeout,
		flushPeriod:    defaultFlushPeriod,
		handlerTimeout: defaultHandlerTimeout,
		batchSize:      defaultBatchSize,
		workersCount:   defaultWorkersCount,
	}

	for _, opt := range opts {
		opt(&o)
	}

	return &MessageCollector{
		caption:        o.captionPrefix + o.caption,
		readyTimeout:   o.readyTimeout,
		flushPeriod:    o.flushPeriod,
		handlerTimeout: o.handlerTimeout,
		batchSize:      o.batchSize,
		workersCount:   o.workersCount,

		handler:         handler,
		errorHandler:    errorHandler,
		logger:          logger,
		contextEmbedder: idGenerator,

		wgMain:       sync.WaitGroup{},
		messageQueue: make(chan []byte),
		workersQueue: make(chan func(ctx context.Context)),
		done:         make(chan struct{}),
	}
}

// Caption - возвращает название сервиса обработки сообщений.
func (p *MessageCollector) Caption() string {
	return p.caption
}

// ReadyTimeout - возвращает максимальное время, за которое должен быть запущен сервис.
func (p *MessageCollector) ReadyTimeout() time.Duration {
	return p.readyTimeout
}

// Start - запуск сервиса обработки сообщений.
// Повторный запуск метода одно и того же объекта не предусмотрен, даже после вызова Shutdown.
func (p *MessageCollector) Start(ctx context.Context, ready func()) error {
	p.wgMain.Add(1)
	defer p.wgMain.Done()

	// WARNING: используется новый контекст со скопированными ID процессами из основного контекста,
	// для того чтобы можно было останавливать процессор только через его метод Shutdown().
	// При вызове этого метода гарантируется корректное завершение работы воркеров процессора.
	ctx = p.contextEmbedder.NewContextWithIDs(ctx)

	p.logger.Debug(ctx, "Starting the message collector...", "collector_name", p.caption)
	defer p.logger.Debug(ctx, "The message collector has been stopped")

	wg := sync.WaitGroup{}
	workersStopped := make(chan struct{})
	ticker := time.NewTicker(p.flushPeriod)

	p.startWorkers(ctx, &wg)

	go func() {
		wg.Wait()
		close(workersStopped)
	}()

	defer func() {
		ticker.Stop()
		close(p.workersQueue)
		<-workersStopped
	}()

	messageBatch := make([][]byte, 0, p.batchSize)

	if ready != nil {
		ready()
	}

	for {
		select {
		case <-p.done:
			if len(messageBatch) == 0 {
				return nil
			}
		case <-ticker.C:
			if len(messageBatch) == 0 {
				continue
			}
		case message := <-p.messageQueue:
			messageBatch = append(messageBatch, message)

			if uint64(len(messageBatch)) < p.batchSize {
				continue
			}
		}

		p.logger.Info(ctx, "Got message batch in message collector...", "messageBatch", len(messageBatch))

		select {
		case <-workersStopped:
			return errInternalWorkersAreStopped.New("collector_name", p.caption)
		case p.workersQueue <- p.workerFunc(messageBatch):
			messageBatch = messageBatch[:0]
		}
	}
}

// PushMessage - comment method.
func (p *MessageCollector) PushMessage(_ context.Context, message []byte) error {
	if p.isSendStopped.Load() {
		return errInternalMessageReceptionStopped.New("collector_name", p.caption)
	}

	p.messageQueue <- message

	return nil
}

// Shutdown - корректная остановка сервиса обработки сообщений.
func (p *MessageCollector) Shutdown(ctx context.Context) error {
	p.logger.Info(ctx, "Shutting down the message collector...")
	p.isSendStopped.Store(true)
	close(p.done)

	p.wgMain.Wait()
	p.logger.Info(ctx, "The message collector has been shut down")

	return nil
}

func (p *MessageCollector) startWorkers(ctx context.Context, wg *sync.WaitGroup) {
	for i := uint64(0); i < p.workersCount; i++ {
		wg.Add(1)

		go func(ctx context.Context) {
			ctx = p.contextEmbedder.WithWorkerIDContext(ctx)

			defer func() {
				wg.Done()

				if rvr := recover(); rvr != nil {
					p.errorHandler.Handle(
						ctx,
						mr.ErrInternalCaughtPanic.New(
							"message collector: "+p.caption,
							rvr,
							string(debug.Stack()),
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

func (p *MessageCollector) workerFunc(messages [][]byte) func(ctx context.Context) {
	return func(ctx context.Context) {
		handlerCtx, cancel := context.WithTimeout(p.contextEmbedder.WithTaskIDContext(ctx), p.handlerTimeout)
		defer cancel()

		if err := p.handler.Execute(handlerCtx, messages); err != nil {
			p.errorHandler.Handle(ctx, err)

			return
		}

		p.logger.Debug(ctx, "The handler has been successfully executed")
	}
}
