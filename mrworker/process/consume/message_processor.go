package consume

import (
	"context"
	"fmt"
	"runtime/debug"
	"strconv"
	"sync"
	"time"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrcore/mrapp"
	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrworker"
)

const (
	defaultCaption           = "MessageProcessor"
	defaultReadyTimeout      = 60 * time.Second
	defaultStartReadDelay    = 0 * time.Second
	defaultReadPeriod        = 60 * time.Second
	defaultCancelReadTimeout = 5 * time.Second
	defaultHandlerTimeout    = 30 * time.Second
	defaultQueueSize         = 100
	defaultWorkersCount      = 1
)

type (
	// MessageProcessor - многопоточный сервис обработки сообщений на основе консьюмера и обработчика.
	MessageProcessor struct {
		caption           string
		readyTimeout      time.Duration
		startReadDelay    time.Duration
		readPeriod        time.Duration
		cancelReadTimeout time.Duration
		handlerTimeout    time.Duration
		queueSize         uint32
		workersCount      uint16
		consumer          mrworker.MessageConsumer
		handler           mrworker.MessageHandler
		errorHandler      mrcore.ErrorHandler
		queue             chan func(ctx context.Context)
		done              chan struct{}
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
		caption:           defaultCaption,
		readyTimeout:      defaultReadyTimeout,
		startReadDelay:    defaultStartReadDelay,
		readPeriod:        defaultReadPeriod,
		cancelReadTimeout: defaultCancelReadTimeout,
		handlerTimeout:    defaultHandlerTimeout,
		queueSize:         defaultQueueSize,
		workersCount:      defaultWorkersCount,
		consumer:          consumer,
		handler:           handler,
		errorHandler:      errorHandler,
		done:              make(chan struct{}, 1),
	}

	for _, opt := range opts {
		opt(p)
	}

	p.queue = make(chan func(ctx context.Context))

	return p
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
func (p *MessageProcessor) Start(ctx context.Context, ready func()) error {
	logger := mrlog.Ctx(ctx)

	wg := sync.WaitGroup{}
	workersStopped := make(chan struct{})

	p.startWorkers(ctx, &wg)

	go func() {
		wg.Wait()
		close(workersStopped)
	}()

	if ready != nil {
		ready()
	}

	// если начало обработки данных необходимо отложить
	if p.startReadDelay > 0 {
		logger.Debug().Msg("The message processor is waiting for the start to read...")
		time.Sleep(p.startReadDelay)
		logger.Debug().Msg("The message processor starts to read...")
	}

	ticker := time.NewTicker(p.readPeriod)

	defer func() {
		ticker.Stop()
		close(p.queue)
		<-workersStopped

		logger.Info().Msg("The message processor has been stopped")
	}()

	for {
		select {
		case <-p.done:
			return nil
		case <-ticker.C:
			messages, err := p.consumer.ReadMessages(ctx, p.queueSize)
			if err != nil {
				if mrcore.ErrInternalTimeoutPeriodHasExpired.Is(err) || mrcore.ErrInternalUnexpectedEOF.Is(err) {
					p.errorHandler.Perform(ctx, err)

					continue
				}

				return err
			}

			logger.Info().Msgf("Got messages %d in message processor...", len(messages))

			for i, message := range messages {
				select {
				case <-workersStopped:
					return func() error {
						cancelCtx, cancel := context.WithTimeout(logger.WithContext(context.Background()), p.cancelReadTimeout)
						defer cancel()

						// передаётся отдельный контекст с персональным таймаутом для исключения внешнего воздействия
						if err = p.consumer.CancelMessages(cancelCtx, messages[i:]); err != nil {
							p.errorHandler.Perform(ctx, err)
						}

						return fmt.Errorf("interrupt the message processor %s, workers are stopped", p.caption)
					}()
				case p.queue <- p.workerFunc(message):
				}
			}
		}
	}
}

// Shutdown - корректная остановка сервиса обработки сообщений.
func (p *MessageProcessor) Shutdown(ctx context.Context) error {
	mrlog.Ctx(ctx).Info().Msg("Shutting down the message processor...")
	close(p.done)

	return nil
}

func (p *MessageProcessor) startWorkers(ctx context.Context, wg *sync.WaitGroup) {
	for i := 0; i < int(p.workersCount); i++ {
		wg.Add(1)

		go func(ctx context.Context, workerNumber int) {
			defer func() {
				wg.Done()

				if rvr := recover(); rvr != nil {
					p.errorHandler.Perform(
						ctx,
						mrcore.ErrInternalCaughtPanic.New(
							"message processor: "+p.caption,
							rvr,
							string(debug.Stack()),
						),
					)
				}
			}()

			logger := mrlog.Ctx(ctx).With().Int(mrapp.KeyWorkerNumber, workerNumber).Logger()
			workerID := mrapp.ProcessCtx(ctx) + mrapp.KeySeparator + "worker-" + strconv.FormatUint(uint64(workerNumber), 10)
			ctx = mrlog.WithContext(mrapp.WithProcessContext(ctx, workerID), logger)

			for fn := range p.queue {
				fn(ctx)
			}

			logger.Debug().Msg("The worker has been stopped")
		}(ctx, i+1)
	}
}

func (p *MessageProcessor) workerFunc(message any) func(ctx context.Context) {
	return func(ctx context.Context) {
		ctx, cancel := context.WithTimeout(ctx, p.handlerTimeout)
		defer cancel()

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

			if err = p.consumer.RejectMessage(ctx, message, err); err != nil {
				p.errorHandler.Perform(ctx, err)
			}

			return
		}

		mrlog.Ctx(ctx).Debug().Msg("The handler has been successfully executed")
	}
}
