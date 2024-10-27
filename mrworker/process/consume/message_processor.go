package consume

import (
	"context"
	"runtime/debug"
	"strconv"
	"sync"
	"time"

	"github.com/rs/xid"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrcore/mrapp"
	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrworker"
)

const (
	defaultCaption      = "MessageProcessor"
	defaultPeriod       = 60 * time.Second
	defaultTimeout      = 60 * time.Second
	defaultQueueSize    = 100
	defaultWorkersCount = 1
)

type (
	// MessageProcessor - многопоточный сервис обработки сообщений на основе консьюмера и обработчика.
	MessageProcessor struct {
		caption      string
		period       time.Duration
		timeout      time.Duration
		queueSize    uint32
		workersCount uint16
		consumer     mrworker.MessageConsumer
		handler      mrworker.MessageHandler
		errorHandler mrcore.ErrorHandler
		queue        chan func(ctx context.Context)
		done         chan struct{}
	}
)

// NewMessageProcessor - создаёт объект MessageProcessor.
func NewMessageProcessor(
	consumer mrworker.MessageConsumer,
	handler mrworker.MessageHandler,
	errorHandler mrcore.ErrorHandler,
	opts ...Option,
) *MessageProcessor {
	p := &MessageProcessor{
		caption:      defaultCaption,
		period:       defaultPeriod,
		timeout:      defaultTimeout,
		queueSize:    defaultQueueSize,
		workersCount: defaultWorkersCount,
		consumer:     consumer,
		handler:      handler,
		errorHandler: errorHandler,
		done:         make(chan struct{}),
	}

	for _, opt := range opts {
		opt(p)
	}

	p.queue = make(chan func(ctx context.Context), p.queueSize)

	return p
}

// Caption - возвращает название сервиса обработки сообщений.
func (p *MessageProcessor) Caption() string {
	return p.caption
}

// Start - запуск сервиса обработки сообщений.
func (p *MessageProcessor) Start(ctx context.Context, ready func()) error {
	processID := xid.New().String()
	logger := mrlog.Ctx(ctx).With().Str("process", p.caption+"-"+processID).Logger()
	ctx = mrlog.WithContext(mrapp.WithProcessContext(ctx, processID), logger)

	wg := sync.WaitGroup{}

	p.startWorkers(ctx, &wg)

	ticker := time.NewTicker(p.period)

	defer func() {
		close(p.queue)
		ticker.Stop()
		wg.Wait()

		logger.Info().Msg("The message processor has been stopped")
	}()

	if ready != nil {
		ready()
	}

	for {
		select {
		case <-p.done:
			logger.Info().Msg("Shutting down the message processor...")

			return nil
		case <-ctx.Done():
			logger.Info().Msg("Interrupt the message processor...")

			return nil
		case <-ticker.C:
			messages, err := p.consumer.ReadMessages(ctx, p.queueSize)
			if err != nil {
				if mrcore.ErrInternalProcessIsStoppedByTimeout.Is(err) || mrcore.ErrInternalUnexpectedEOF.Is(err) {
					p.errorHandler.Perform(ctx, err)

					continue
				}

				return err
			}

			logger.Info().Msgf("Got messages %d in message processor...", len(messages))

			for _, message := range messages {
				p.queue <- p.workerFunc(message)
			}
		}
	}
}

// Shutdown - корректная остановка сервиса обработки сообщений.
func (p *MessageProcessor) Shutdown(_ context.Context) error {
	close(p.done)

	return nil
}

func (p *MessageProcessor) startWorkers(ctx context.Context, wg *sync.WaitGroup) {
	for i := 0; i < int(p.workersCount); i++ {
		wg.Add(1)

		go func(ctx context.Context, workerNumber int) {
			workerID := mrapp.ProcessCtx(ctx) + "-worker" + strconv.Itoa(workerNumber)
			logger := mrlog.Ctx(ctx).With().Int("worker", workerNumber).Logger()
			ctx = mrlog.WithContext(mrapp.WithProcessContext(ctx, workerID), logger)

			defer wg.Done()

			for fn := range p.queue {
				fn(ctx)
			}

			logger.Debug().Msg("Finished the worker")
		}(ctx, i+1)
	}
}

func (p *MessageProcessor) workerFunc(message any) func(ctx context.Context) {
	return func(ctx context.Context) {
		ctx, cancel := context.WithTimeout(ctx, p.timeout)

		defer func() {
			cancel()

			if rvr := recover(); rvr != nil {
				p.errorHandler.Perform(
					ctx,
					mrcore.ErrInternalCaughtPanic.New(
						"message processor: "+p.caption,
						rvr,
						debug.Stack(),
					),
				)
			}
		}()

		commit, err := p.handler.Execute(ctx, message)
		if err != nil {
			p.errorHandler.Perform(ctx, err)

			if err = p.consumer.RejectMessage(ctx, message, err); err != nil {
				p.errorHandler.Perform(ctx, err)
			}

			return
		}

		// если консьюмер и обработчик работают в рамках одной БД,
		// то коммит обработчика с коммитом консьюмера могут проходить в единой транзакции
		if err = p.consumer.CommitMessage(ctx, message, commit); err != nil {
			p.errorHandler.Perform(ctx, err)

			return
		}

		mrlog.Ctx(ctx).Info().Msgf("Handler has executed")
	}
}
