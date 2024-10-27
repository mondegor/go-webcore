package schedule

import (
	"context"
	"fmt"
	"runtime/debug"
	"sync"
	"time"

	"github.com/rs/xid"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrcore/mrapp"
	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrworker"
)

const (
	defaultCaption = "TaskScheduler"
)

type (
	// TaskScheduler - многопоточный сервис запуска задач по расписанию (планировщик задач).
	TaskScheduler struct {
		caption      string
		tasks        []mrworker.Task
		errorHandler mrcore.ErrorHandler
		done         chan struct{}
	}
)

// NewTaskScheduler - создаёт объект TaskScheduler.
func NewTaskScheduler(errorHandler mrcore.ErrorHandler, opts ...Option) *TaskScheduler {
	s := &TaskScheduler{
		caption:      defaultCaption,
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

// Start - запуск планировщика задач.
func (s *TaskScheduler) Start(ctx context.Context, ready func()) error {
	processID := xid.New().String()
	logger := mrlog.Ctx(ctx).With().Str("process", s.caption+"-"+processID).Logger()
	ctx = mrlog.WithContext(mrapp.WithProcessContext(ctx, processID), logger)

	if err := s.startup(ctx); err != nil {
		return err
	}

	tasksDone := make(chan struct{})
	wg := sync.WaitGroup{}

	for i := range s.tasks {
		wg.Add(1)

		go func(ctx context.Context, task mrworker.Task) {
			taskID := mrapp.ProcessCtx(ctx) + "-task-" + task.Caption()
			logger := mrlog.Ctx(ctx).With().Str("task", task.Caption()).Logger()
			ctx = mrlog.WithContext(mrapp.WithProcessContext(ctx, taskID), logger)

			defer wg.Done()

			ticker := time.NewTicker(task.Period())
			defer ticker.Stop()

			for {
				select {
				case <-tasksDone:
					logger.Info().Msgf("Interrupt the task")

					return
				case <-ticker.C:
					func(ctx context.Context) {
						ctx, cancel := context.WithTimeout(ctx, task.Timeout())

						defer func() {
							cancel()

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

	select {
	case <-s.done:
		logger.Info().Msgf("Shutting down the task scheduler...")
	case <-ctx.Done():
		logger.Info().Msgf("Interrupt the task scheduler...")
	}

	close(tasksDone)
	wg.Wait()

	logger.Info().Msgf("The task scheduler has been stopped")

	return nil
}

// Shutdown - корректная остановка планировщика задач.
func (s *TaskScheduler) Shutdown(_ context.Context) error {
	close(s.done)

	return nil
}

func (s *TaskScheduler) startup(ctx context.Context) error {
	if len(s.tasks) == 0 {
		return fmt.Errorf("no tasks to start for the task scheduler %s", s.caption)
	}

	logger := mrlog.Ctx(ctx)

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

		logger.Info().Msgf("Startup the task %s...", task.Caption())

		if err := task.Do(ctx); err != nil {
			return err
		}
	}

	return nil
}
