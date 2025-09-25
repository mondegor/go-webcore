package schedule

import (
	"context"
	"runtime/debug"
	"sync"
	"time"

	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-sysmess/mrerr/mr"
	"github.com/mondegor/go-sysmess/mrlog"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrworker"
)

//go:generate mockgen -source=task_scheduler.go -destination=./mock/task_scheduler.go

const (
	defaultCaption      = "TaskScheduler"
	defaultReadyTimeout = 30 * time.Second
)

type (
	// TaskScheduler - многопоточный сервис запуска задач по расписанию (планировщик задач).
	TaskScheduler struct {
		caption         string
		readyTimeout    time.Duration
		tasks           []mrworker.Task
		errorHandler    mrcore.ErrorHandler
		logger          mrlog.Logger
		contextEmbedder contextEmbedder
		wgMain          sync.WaitGroup
		done            chan struct{}
	}

	contextEmbedder interface {
		WithWorkerIDContext(ctx context.Context) context.Context
		WithTaskIDContext(ctx context.Context) context.Context
	}

	options struct {
		caption       string
		captionPrefix string
		readyTimeout  time.Duration
		tasks         []mrworker.Task
	}
)

var (
	// ErrInternalNoTasks - comment var.
	ErrInternalNoTasks = mrerr.NewKindInternal("no tasks to start for the task scheduler")

	// ErrInternalZeroParam - comment var.
	ErrInternalZeroParam = mrerr.NewKindInternal("task has zero param for the task scheduler: {ParamName}")
)

// NewTaskScheduler - создаёт объект TaskScheduler.
func NewTaskScheduler(errorHandler mrcore.ErrorHandler, logger mrlog.Logger, contextEmbedder contextEmbedder, opts ...Option) *TaskScheduler {
	o := options{
		caption:      defaultCaption,
		readyTimeout: defaultReadyTimeout,
	}

	for _, opt := range opts {
		opt(&o)
	}

	return &TaskScheduler{
		caption:         o.captionPrefix + o.caption,
		readyTimeout:    o.readyTimeout,
		tasks:           o.tasks,
		errorHandler:    errorHandler,
		logger:          logger,
		contextEmbedder: contextEmbedder,
		wgMain:          sync.WaitGroup{},
		done:            make(chan struct{}),
	}
}

// Caption - возвращает название планировщика задач.
func (s *TaskScheduler) Caption() string {
	return s.caption
}

// ReadyTimeout - возвращает максимальное время, за которое должен быть запущен планировщик со всеми его задачами.
func (s *TaskScheduler) ReadyTimeout() time.Duration {
	return s.readyTimeout
}

// Start - запуск планировщика задач.
// Повторный запуск метода одно и того же объекта не предусмотрен, даже после вызова Shutdown.
func (s *TaskScheduler) Start(ctx context.Context, ready func()) error {
	s.wgMain.Add(1)
	defer s.wgMain.Done()

	s.logger.Debug(ctx, "Starting the task scheduler...")
	defer s.logger.Debug(ctx, "The task scheduler has been stopped")

	if err := s.startup(ctx); err != nil {
		return err
	}

	wg := sync.WaitGroup{}

	for i := range s.tasks {
		wg.Add(1)

		go func(ctx context.Context, task mrworker.Task) {
			defer wg.Done()

			ctx = s.contextEmbedder.WithWorkerIDContext(ctx)
			s.logger.Info(ctx, "Starting worker", "task_name", task.Caption())

			ticker := time.NewTicker(task.Period())

			defer func() {
				ticker.Stop()

				if rvr := recover(); rvr != nil {
					s.errorHandler.Handle(
						ctx,
						mr.ErrInternalCaughtPanic.New(
							"task worker: "+task.Caption(),
							rvr,
							string(debug.Stack()),
						),
					)
				}
			}()

			for {
				select {
				case <-s.done:
					s.logger.Debug(ctx, "The task worker has been stopped")

					return
				case <-task.SignalDo():
					s.logger.Debug(ctx, "signalDo event", "task_name", task.Caption())
					ticker.Reset(task.Period())
				case <-ticker.C:
					s.logger.Debug(ctx, "ticker event", "task_name", task.Caption())
				}

				func(ctx context.Context) {
					ctx = s.contextEmbedder.WithTaskIDContext(ctx)
					s.logger.Info(ctx, "Execute task", "task_name", task.Caption())

					ctx, cancel := context.WithTimeout(ctx, task.Timeout())
					defer cancel()

					if err := task.Do(ctx); err != nil {
						s.errorHandler.Handle(ctx, err)
					}
				}(ctx)
			}
		}(ctx, s.tasks[i])
	}

	if ready != nil {
		ready()
	}

	wg.Wait()

	return nil
}

// Shutdown - корректная остановка планировщика задач.
func (s *TaskScheduler) Shutdown(ctx context.Context) error {
	s.logger.Info(ctx, "Shutting down the task scheduler...")
	close(s.done)

	s.wgMain.Wait()
	s.logger.Info(ctx, "The task scheduler has been shut down")

	return nil
}

func (s *TaskScheduler) startup(ctx context.Context) error {
	if len(s.tasks) == 0 {
		return ErrInternalNoTasks.New("scheduler_name", s.caption)
	}

	for _, task := range s.tasks {
		if task.Period() == 0 {
			return ErrInternalZeroParam.New("period", "task_name", task.Caption())
		}

		if task.Timeout() == 0 {
			return ErrInternalZeroParam.New("timeout", "task_name", task.Caption())
		}

		// запуск задач на этапе старта планировщика выполняется последовательно
		// и без использования таймаута, чтобы гарантировать полную их инициализацию
		if !task.Startup() {
			continue
		}

		s.logger.Debug(ctx, "Startup the task...", "task_name", task.Caption())

		if err := task.Do(ctx); err != nil {
			return err
		}

		s.logger.Debug(ctx, "The task is completed", "task_name", task.Caption())
	}

	return nil
}
