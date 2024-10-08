package mrschedule

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrworker"
)

const (
	schedulerName = "Scheduler"
)

type (
	// Scheduler - планировщик задач.
	Scheduler struct {
		done         chan struct{}
		errorHandler mrcore.ErrorHandler
		tasks        []mrworker.Task
	}
)

// NewScheduler - создаёт объект Scheduler.
func NewScheduler(errorHandler mrcore.ErrorHandler, tasks ...mrworker.Task) *Scheduler {
	return &Scheduler{
		done:         make(chan struct{}),
		errorHandler: errorHandler,
		tasks:        tasks,
	}
}

// Caption - возвращает название планировщика.
func (s *Scheduler) Caption() string {
	return schedulerName
}

// Start - запуск планировщика задач.
func (s *Scheduler) Start(ctx context.Context, ready func()) error {
	logger := mrlog.Ctx(ctx).With().Str("process", schedulerName).Logger()
	ctx = mrlog.WithContext(ctx, logger)

	if err := s.startup(ctx, logger); err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(ctx)
	wg := sync.WaitGroup{}

	for i := range s.tasks {
		wg.Add(1)

		go func(ctx context.Context, task mrworker.Task) {
			defer wg.Done()

			ticker := time.NewTicker(task.Period())
			defer ticker.Stop()

			for {
				select {
				case <-ticker.C:
					func(ctx context.Context) {
						ctx, cancel := context.WithTimeout(ctx, task.Timeout())
						defer cancel()

						if err := task.Do(ctx); err != nil {
							s.errorHandler.Perform(ctx, err)
						}
					}(ctx)
				case <-ctx.Done():
					logger.Info().Msgf("Interrupt task '%s'", task.Caption())

					return
				}
			}
		}(ctx, s.tasks[i])
	}

	if ready != nil {
		ready()
	}

	select {
	case <-s.done:
		logger.Info().Msg("Shutting down the scheduler...")
	case <-ctx.Done():
		logger.Info().Msgf("Interrupt the scheduler...")
	}

	cancel()
	wg.Wait()

	logger.Info().Msg("The scheduler has been stopped")

	return nil
}

// Shutdown - корректная остановка планировщика задач.
func (s *Scheduler) Shutdown(_ context.Context) error {
	close(s.done)

	return nil
}

func (s *Scheduler) startup(ctx context.Context, logger mrlog.Logger) error {
	for _, task := range s.tasks {
		if task.Period() == 0 {
			return fmt.Errorf("task '%s' has zero period", task.Caption())
		}

		if task.Timeout() == 0 {
			return fmt.Errorf("task '%s' has zero timeout", task.Caption())
		}

		// последовательный первый запуск задач, если они требуют этого
		// запускаются без таймаута, т.к. требуется, чтобы они полностью завершились перед запуском системы
		if !task.Startup() {
			continue
		}

		logger.Info().Msgf("Startup the task '%s'...", task.Caption())

		if err := task.Do(ctx); err != nil {
			return err
		}
	}

	return nil
}
