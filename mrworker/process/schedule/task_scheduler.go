package schedule

import (
	"context"
	"fmt"
	"runtime/debug"
	"sync"
	"time"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrcore/mrapp"
	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrworker"
)

const (
	defaultCaption      = "TaskScheduler"
	defaultReadyTimeout = 60 * time.Second
)

type (
	// TaskScheduler - многопоточный сервис запуска задач по расписанию (планировщик задач).
	TaskScheduler struct {
		caption      string
		readyTimeout time.Duration
		tasks        []mrworker.Task
		errorHandler mrcore.ErrorHandler
		done         chan struct{}
	}
)

// NewTaskScheduler - создаёт объект TaskScheduler.
func NewTaskScheduler(errorHandler mrcore.ErrorHandler, opts ...Option) *TaskScheduler {
	s := &TaskScheduler{
		caption:      defaultCaption,
		readyTimeout: defaultReadyTimeout,
		errorHandler: errorHandler,
		done:         make(chan struct{}),
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
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
func (s *TaskScheduler) Start(ctx context.Context, ready func()) error {
	mrlog.Ctx(ctx).Info().Msg("Starting the task scheduler...")

	if err := s.startup(ctx); err != nil {
		return err
	}

	wg := sync.WaitGroup{}

	for i := range s.tasks {
		wg.Add(1)

		go func(ctx context.Context, task mrworker.Task) {
			defer wg.Done()

			logger := mrlog.Ctx(ctx).With().Str("task", task.Caption()).Logger()
			taskID := mrapp.ProcessCtx(ctx) + mrapp.KeySeparator + "task-" + task.Caption()
			ctx = mrlog.WithContext(mrapp.WithProcessContext(ctx, taskID), logger)
			ticker := time.NewTicker(task.Period())

			defer func() {
				ticker.Stop()

				if rvr := recover(); rvr != nil {
					s.errorHandler.Perform(
						ctx,
						mrcore.ErrInternalCaughtPanic.New(
							"task: "+taskID,
							rvr,
							debug.Stack(),
						),
					)
				}
			}()

			for {
				select {
				case <-s.done:
					logger.Debug().Msg("The task worker has been stopped")

					return
				case <-ticker.C:
					func(ctx context.Context) {
						ctx, cancel := context.WithTimeout(ctx, task.Timeout())
						defer cancel()

						if err := task.Do(ctx); err != nil {
							s.errorHandler.Perform(ctx, err)
						}
					}(ctx)
				}
			}
		}(ctx, s.tasks[i])
	}

	if ready != nil {
		ready()
	}

	wg.Wait()
	mrlog.Ctx(ctx).Info().Msg("The task scheduler has been stopped")

	return nil
}

// Shutdown - корректная остановка планировщика задач.
func (s *TaskScheduler) Shutdown(ctx context.Context) error {
	mrlog.Ctx(ctx).Info().Msg("Shutting down the task scheduler...")
	close(s.done)

	return nil
}

func (s *TaskScheduler) startup(ctx context.Context) error {
	if len(s.tasks) == 0 {
		return fmt.Errorf("no tasks to start for the task scheduler %s", s.caption)
	}

	for _, task := range s.tasks {
		if task.Period() == 0 {
			return fmt.Errorf("task %s has zero period for the task scheduler %s", task.Caption(), s.caption)
		}

		if task.Timeout() == 0 {
			return fmt.Errorf("task %s has zero timeout for the task scheduler %s", task.Caption(), s.caption)
		}

		// последовательный первый запуск задач, если они этого требуют
		// запуск происходит без таймаута, т.к. требуется, чтобы задачи полностью завершились перед запуском системы
		if !task.Startup() {
			continue
		}

		mrlog.Ctx(ctx).Debug().Msgf("Startup the task %s...", task.Caption())

		if err := task.Do(ctx); err != nil {
			return err
		}

		mrlog.Ctx(ctx).Debug().Msgf("The task %s is completed", task.Caption())
	}

	return nil
}
