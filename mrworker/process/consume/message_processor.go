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
	// defaultCaption - название сервиса по умолчанию.
	defaultCaption = "MessageProcessor"

	// defaultReadyTimeout - таймаут готовности сервиса по умолчанию.
	defaultReadyTimeout = 30 * time.Second

	// defaultConsumerReadTimeout - таймаут чтения консьюмером по умолчанию.
	defaultConsumerReadTimeout = 2 * time.Second

	// defaultConsumerWriteTimeout - таймаут записи консьюмером по умолчанию.
	defaultConsumerWriteTimeout = 3 * time.Second

	// defaultHandlerTimeout - таймаут выполнения обработчика сообщения по умолчанию.
	defaultHandlerTimeout = 30 * time.Second

	// defaultQueueSize - размер очереди сообщений для одной выборки по умолчанию.
	defaultQueueSize = 100

	// defaultWorkersCount - количество воркеров-обработчиков по умолчанию.
	defaultWorkersCount = 1
)

var (
	// defaultReadPeriodStrategy - период опроса очереди в состоянии простоя.
	defaultReadPeriodStrategy = mrworker.NewStaticPeriodStrategy(60 * time.Second) //nolint:gochecknoglobals

	errInternalWorkersAreStopped = errors.NewInternalProto("the message processor workers has been stopped")
)

type (
	// MessageProcessor - многопоточный сервис обработки сообщений (PULL-модель).
	//
	// Принцип работы:
	//  1. Периодически опрашивает очередь через MessageConsumer (PULL);
	//  2. Каждое сообщение отправляется в workersQueue для обработки;
	//  3. Обработчик (MessageHandler) выполняет работу и возвращает функцию commit;
	//  4. При успехе вызывает CommitMessage (возможно с preCommit), при ошибке - RejectMessage;
	MessageProcessor[T any] struct {
		caption            string
		readyTimeout       time.Duration
		readPeriodStrategy mrworker.PeriodStrategy
		handlerTimeout     time.Duration
		queueSize          int
		workersCount       int

		consumer     mrworker.MessageConsumer[T]
		handler      mrworker.MessageHandler[T]
		errorHandler errors.Handler
		logger       mrlog.Logger
		traceManager mrtrace.ContextManager

		wg            sync.WaitGroup
		signalExecute <-chan struct{}
		workersQueue  chan func(ctx context.Context)
		done          chan struct{}
	}
)

// NewMessageProcessor - создаёт сервис обработки сообщений (PULL-модель).
func NewMessageProcessor[T any](
	consumer mrworker.MessageConsumer[T],
	handler mrworker.MessageHandler[T],
	errorHandler errors.Handler,
	logger mrlog.Logger,
	traceManager mrtrace.ContextManager,
	opts ...Option[T],
) *MessageProcessor[T] {
	o := options[T]{
		processor: &MessageProcessor[T]{
			caption:            defaultCaption,
			readyTimeout:       defaultReadyTimeout,
			readPeriodStrategy: defaultReadPeriodStrategy,
			handlerTimeout:     defaultHandlerTimeout,

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

// Caption - возвращает название сервиса обработки сообщений в свободной форме.
func (p *MessageProcessor[T]) Caption() string {
	return p.caption
}

// ReadyTimeout - возвращает максимальное время, за которое должен быть запущен сервис.
func (p *MessageProcessor[T]) ReadyTimeout() time.Duration {
	return p.readyTimeout
}

// Start - запуск сервиса обработки сообщений.
//
// Процесс работы:
//  1. Запускает N воркеров для обработки сообщений;
//  2. Периодически (readPeriodStrategy) опрашивает очередь через consumer.ReadMessages;
//  3. Каждое сообщение отправляется в workersQueue для обработки;
//  4. При signalExecute немедленно выполняет опрос очереди;
//  5. При ошибке таймаута/EOF логирует ошибку и продолжает работу;
//  6. При других ошибках завершает работу;
//
// Важно:
//   - Отмена внешнего контекста приведёт к завершению процесса;
//   - Для корректной остановки используйте Shutdown;
//   - Повторный запуск того же объекта не поддерживается.
func (p *MessageProcessor[T]) Start(ctx context.Context, ready func()) error {
	p.wg.Add(1)
	defer p.wg.Done()

	p.logger.Debug(ctx, "Starting the message processor...", "processor_name", p.caption)
	defer p.logger.Debug(ctx, "The message processor has been stopped")

	wgWorkers := sync.WaitGroup{}
	workersStopped := make(chan struct{})
	ticker := time.NewTicker(p.readPeriodStrategy.Period())

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
		case <-ticker.C:
			p.logger.Debug(ctx, "ticker.C event", "processor_name", p.caption)
		}

		ticker.Reset(p.readPeriodStrategy.Period())

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
// Останавливает основной цикл и ожидает завершения всех воркеров.
//
// Важно: при повторном вызове произойдёт panic (закрытие закрытого канала done).
func (p *MessageProcessor[T]) Shutdown(ctx context.Context) error {
	p.logger.Debug(ctx, "Shutting down the message processor...")
	close(p.done)

	p.wg.Wait()
	p.logger.Debug(ctx, "The message processor has been shut down")

	return nil
}

func (p *MessageProcessor[T]) startWorkers(ctx context.Context, wg *sync.WaitGroup) {
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

// workerFunc - создаёт функцию обработки одного сообщения для отправки в воркер.
//
// Логика работы:
//  1. Вызывает handler.Execute(message) для обработки сообщения;
//  2. При ошибке: логирует ошибку и вызывает RejectMessage;
//  3. При успехе: вызывает CommitMessage с функцией commit от обработчика;
//     (commit и consumer commit могут выполняться в единой транзакции);
//  4. При ошибке коммита: вызывает RejectMessage.
func (p *MessageProcessor[T]) workerFunc(message T) func(ctx context.Context) {
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
