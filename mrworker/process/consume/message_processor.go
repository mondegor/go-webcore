package consume

import (
	"context"
	"errors"
	"runtime/debug"
	"sync"
	"time"

	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-sysmess/mrerr/mr"
	"github.com/mondegor/go-sysmess/mrlog"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrworker"
)

const (
	defaultCaption              = "MessageProcessor"
	defaultReadyTimeout         = 30 * time.Second
	defaultReadPeriod           = 60 * time.Second
	defaultConsumerReadTimeout  = 2 * time.Second
	defaultConsumerWriteTimeout = 3 * time.Second
	defaultHandlerTimeout       = 30 * time.Second
	defaultQueueSize            = 100
	defaultWorkersCount         = 1
)

type (
	// MessageProcessor - многопоточный сервис обработки сообщений на основе консьюмера и обработчика.
	MessageProcessor struct {
		caption        string
		readyTimeout   time.Duration
		readPeriod     time.Duration
		handlerTimeout time.Duration
		queueSize      uint64
		workersCount   uint64

		consumer        mrworker.MessageConsumer
		handler         mrworker.MessageHandler
		errorHandler    mrcore.ErrorHandler
		logger          mrlog.Logger
		contextEmbedder contextEmbedder

		wgMain        sync.WaitGroup
		signalExecute <-chan struct{}
		workersQueue  chan func(ctx context.Context)
		done          chan struct{}
	}

	contextEmbedder interface {
		NewContextWithIDs(originalCtx context.Context) context.Context
		WithWorkerIDContext(ctx context.Context) context.Context
		WithTaskIDContext(ctx context.Context) context.Context
	}

	options struct {
		caption              string
		captionPrefix        string
		readyTimeout         time.Duration
		readPeriod           time.Duration
		consumerReadTimeout  time.Duration
		consumerWriteTimeout time.Duration
		handlerTimeout       time.Duration
		queueSize            uint64
		workersCount         uint64
		signalExecute        <-chan struct{}
	}
)

var errInternalWorkersAreStopped = mrerr.NewKindInternal("the message processor workers has been stopped")

// NewMessageProcessor - создаёт объект MessageProcessor.
func NewMessageProcessor(
	consumer mrworker.MessageConsumer,
	handler mrworker.MessageHandler,
	errorHandler mrcore.ErrorHandler,
	logger mrlog.Logger,
	idGenerator contextEmbedder,
	opts ...Option,
) *MessageProcessor {
	o := options{
		caption:              defaultCaption,
		readyTimeout:         defaultReadyTimeout,
		readPeriod:           defaultReadPeriod,
		consumerReadTimeout:  defaultConsumerReadTimeout,
		consumerWriteTimeout: defaultConsumerWriteTimeout,
		handlerTimeout:       defaultHandlerTimeout,
		queueSize:            defaultQueueSize,
		workersCount:         defaultWorkersCount,
	}

	for _, opt := range opts {
		opt(&o)
	}

	return &MessageProcessor{
		caption:        o.captionPrefix + o.caption,
		readyTimeout:   o.readyTimeout,
		readPeriod:     o.readPeriod,
		handlerTimeout: o.handlerTimeout,
		queueSize:      o.queueSize,
		workersCount:   o.workersCount,

		consumer:        NewConsumerWithTimeout(consumer, o.consumerReadTimeout, o.consumerWriteTimeout),
		handler:         handler,
		errorHandler:    errorHandler,
		logger:          logger,
		contextEmbedder: idGenerator,

		wgMain:        sync.WaitGroup{},
		signalExecute: o.signalExecute,
		workersQueue:  make(chan func(ctx context.Context)),
		done:          make(chan struct{}),
	}
}

// Caption - возвращает название сервиса обработки сообщений.
func (p *MessageProcessor) Caption() string {
	return p.caption
}

// ReadyTimeout - возвращает максимальное время, за которое должен быть запущен сервис.
func (p *MessageProcessor) ReadyTimeout() time.Duration {
	return p.readyTimeout
}

// Start - запуск сервиса обработки сообщений.
// Повторный запуск метода одно и того же объекта не предусмотрен, даже после вызова Shutdown.
func (p *MessageProcessor) Start(ctx context.Context, ready func()) error {
	p.wgMain.Add(1)
	defer p.wgMain.Done()

	// WARNING: используется новый контекст со скопированными ID процессами из основного контекста,
	// для того чтобы можно было останавливать процессор только через его метод Shutdown().
	// При вызове этого метода гарантируется корректное завершение работы воркеров процессора.
	ctx = p.contextEmbedder.NewContextWithIDs(ctx)

	p.logger.Debug(ctx, "Starting the message processor...", "processor_name", p.caption)
	defer p.logger.Debug(ctx, "The message processor has been stopped")

	wg := sync.WaitGroup{}
	workersStopped := make(chan struct{})
	ticker := time.NewTicker(p.readPeriod)

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

	if ready != nil {
		ready()
	}

	for {
		select {
		case <-p.done:
			return nil
		case <-p.signalExecute:
			p.logger.Debug(ctx, "signalExecute event", "processor_name", p.caption)
			ticker.Reset(p.readPeriod)
		case <-ticker.C:
			p.logger.Debug(ctx, "ticker event", "processor_name", p.caption)
		}

		ctx = p.contextEmbedder.WithTaskIDContext(ctx) // producerID

		messages, err := p.consumer.ReadMessages(ctx, p.queueSize)
		if err != nil {
			if errors.Is(err, mr.ErrInternalTimeoutPeriodHasExpired) || errors.Is(err, mr.ErrInternalUnexpectedEOF) {
				p.errorHandler.Handle(ctx, err)

				continue
			}

			return err
		}

		p.logger.Info(ctx, "Got messages in message processor...", "countMessages", len(messages))

		for i, message := range messages {
			select {
			case <-workersStopped:
				return func() error {
					if err = p.consumer.CancelMessages(ctx, messages[i:]); err != nil {
						p.errorHandler.Handle(ctx, err)
					}

					return errInternalWorkersAreStopped.New("processor_name", p.caption)
				}()
			case p.workersQueue <- p.workerFunc(message):
			}
		}
	}
}

// Shutdown - корректная остановка сервиса обработки сообщений.
func (p *MessageProcessor) Shutdown(ctx context.Context) error {
	p.logger.Info(ctx, "Shutting down the message processor...")
	close(p.done)

	p.wgMain.Wait()
	p.logger.Info(ctx, "The message processor has been shut down")

	return nil
}

func (p *MessageProcessor) startWorkers(ctx context.Context, wg *sync.WaitGroup) {
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
							"message processor: "+p.caption,
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

func (p *MessageProcessor) workerFunc(message any) func(ctx context.Context) {
	return func(ctx context.Context) {
		ctx = p.contextEmbedder.WithTaskIDContext(ctx)

		handlerCtx, cancel := context.WithTimeout(ctx, p.handlerTimeout)
		defer cancel()

		commit, err := p.handler.Execute(handlerCtx, message)
		if err != nil {
			p.errorHandler.Handle(ctx, err)

			if err = p.consumer.RejectMessage(ctx, message, err); err != nil {
				p.errorHandler.Handle(ctx, err)
			}

			return
		}

		// если консьюмер и обработчик работают в рамках одной БД,
		// то коммит обработчика с коммитом консьюмера могут проходить в единой транзакции
		if err = p.consumer.CommitMessage(ctx, message, commit); err != nil {
			p.errorHandler.Handle(ctx, err)

			if err = p.consumer.RejectMessage(ctx, message, err); err != nil {
				p.errorHandler.Handle(ctx, err)
			}

			return
		}

		p.logger.Debug(ctx, "The handler has been successfully executed")
	}
}
