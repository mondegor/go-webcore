package consume

import (
	"context"
	"runtime/debug"
	"sync"
	"time"

	"github.com/mondegor/go-sysmess/errors"
	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-sysmess/mrtrace"

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
		queueSize      int
		workersCount   int

		consumer     mrworker.MessageConsumer
		handler      mrworker.MessageHandler
		errorHandler errors.Handler
		logger       mrlog.Logger
		traceManager mrtrace.ContextManager

		wg            sync.WaitGroup
		signalExecute <-chan struct{}
		workersQueue  chan func(ctx context.Context)
		done          chan struct{}
	}
)

var errInternalWorkersAreStopped = errors.NewInternalProto("the message processor workers has been stopped")

// NewMessageProcessor - создаёт объект MessageProcessor.
func NewMessageProcessor(
	consumer mrworker.MessageConsumer,
	handler mrworker.MessageHandler,
	errorHandler errors.Handler,
	logger mrlog.Logger,
	traceManager mrtrace.ContextManager,
	opts ...Option,
) *MessageProcessor {
	o := options{
		processor: &MessageProcessor{
			caption:        defaultCaption,
			readyTimeout:   defaultReadyTimeout,
			readPeriod:     defaultReadPeriod,
			handlerTimeout: defaultHandlerTimeout,

			handler:      handler,
			errorHandler: errorHandler,
			logger:       logger,
			traceManager: traceManager,

			wg:           sync.WaitGroup{},
			workersQueue: make(chan func(ctx context.Context)),
			done:         make(chan struct{}),
		},
		consumerReadTimeout:  defaultConsumerReadTimeout,
		consumerWriteTimeout: defaultConsumerWriteTimeout,
	}

	for _, opt := range opts {
		opt(&o)
	}

	if o.captionPrefix != "" {
		o.processor.caption = o.captionPrefix + o.processor.caption
	}

	if o.processor.queueSize < 1 {
		o.processor.queueSize = defaultQueueSize
	}

	if o.processor.workersCount < 1 {
		o.processor.workersCount = defaultWorkersCount
	}

	if o.consumerReadTimeout > 0 || o.consumerWriteTimeout > 0 {
		o.processor.consumer = NewConsumerWithTimeout(
			consumer,
			o.consumerReadTimeout,
			o.consumerWriteTimeout,
		)
	}

	return o.processor
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
// Отмена внешнего контекста приведёт к аварийному завершению процесса,
// для корректной остановки следует использовать Shutdown.
// Повторный запуск метода одно и того же объекта не предусмотрен, даже после вызова Shutdown.
func (p *MessageProcessor) Start(ctx context.Context, ready func()) error {
	p.wg.Add(1)
	defer p.wg.Done()

	p.logger.Debug(ctx, "Starting the message processor...", "processor_name", p.caption)
	defer p.logger.Debug(ctx, "The message processor has been stopped")

	wgWorkers := sync.WaitGroup{}
	workersStopped := make(chan struct{})
	ticker := time.NewTicker(p.readPeriod)

	p.startWorkers(ctx, &wgWorkers)

	go func() {
		wgWorkers.Wait()
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
		case <-ctx.Done():
			p.logger.Debug(ctx, "The message processor detected context 'Done'", "error", ctx.Err())

			return nil
		case <-p.signalExecute:
			p.logger.Debug(ctx, "signalExecute event", "processor_name", p.caption)
			ticker.Reset(p.readPeriod)
		case <-ticker.C:
			p.logger.Debug(ctx, "ticker.C event", "processor_name", p.caption)
		}

		ctx = p.traceManager.WithGeneratedProcessID(ctx, mrtrace.KeyTaskID) // producerID

		messages, err := p.consumer.ReadMessages(ctx, p.queueSize)
		if err != nil {
			if errors.Is(err, errors.ErrSystemTimeoutPeriodHasExpired) || errors.Is(err, errors.ErrSystemStorageUnexpectedEOF) {
				p.errorHandler.Handle(ctx, err)

				continue
			}

			return err
		}

		p.logger.Info(ctx, "Got messages in the message processor...", "count_messages", len(messages))

		for i, message := range messages {
			select {
			case <-workersStopped:
				if err = p.consumer.CancelMessages(ctx, messages[i:]); err != nil {
					p.errorHandler.Handle(ctx, err)
				}

				return errInternalWorkersAreStopped.New("processor_name", p.caption)
			case p.workersQueue <- p.workerFunc(message):
			}
		}
	}
}

// Shutdown - корректная остановка сервиса обработки сообщений.
// При повторном вызове метода произойдёт panic.
func (p *MessageProcessor) Shutdown(ctx context.Context) error {
	p.logger.Debug(ctx, "Shutting down the message processor...")
	close(p.done)

	p.wg.Wait()
	p.logger.Debug(ctx, "The message processor has been shut down")

	return nil
}

func (p *MessageProcessor) startWorkers(ctx context.Context, wg *sync.WaitGroup) {
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
							"source", "message processor: "+p.caption,
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

func (p *MessageProcessor) workerFunc(message any) func(ctx context.Context) {
	return func(ctx context.Context) {
		ctx = p.traceManager.WithGeneratedProcessID(ctx, mrtrace.KeyTaskID)

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
