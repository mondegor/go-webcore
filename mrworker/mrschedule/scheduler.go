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

// PrepareToStart - comment method.
func (s *Scheduler) PrepareToStart(ctx context.Context) (execute func() error, interrupt func(error)) {
	return func() error {
			return s.Start(ctx)
		}, func(_ error) {
			s.Stop()
		}
}

// Start - запуск планировщика задач.
func (s *Scheduler) Start(ctx context.Context) error {
	logger := mrlog.Ctx(ctx).With().Str("service", schedulerName).Logger()
	ctx = mrlog.WithContext(ctx, logger)

	if err := s.startup(ctx); err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(ctx)
	wg := sync.WaitGroup{}

	for i := range s.tasks {
		task := s.tasks[i]

		wg.Add(1)

		go func(ctx context.Context) {
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
							s.errorHandler.Process(ctx, err)
						}
					}(ctx)
				case <-ctx.Done():
					logger.Info().Msgf("Interrupt task %s", task.Caption())

					return
				}
			}
		}(ctx)
	}

	select {
	case <-s.done:
		logger.Info().Msg("Shutting down the scheduler...")
	case <-ctx.Done():
		logger.Info().Msgf("Interrupt the scheduler...")
	}

	cancel()
	wg.Wait()

	logger.Info().Msg("The scheduler has been shutdown")

	return nil
}

// Stop - остановка планировщика задач.
func (s *Scheduler) Stop() {
	close(s.done)
}

func (s *Scheduler) startup(ctx context.Context) error {
	for _, task := range s.tasks {
		if task.Period() == 0 {
			return fmt.Errorf("task %s has zero period", task.Caption())
		}

		if task.Timeout() == 0 {
			return fmt.Errorf("task %s has zero timeout", task.Caption())
		}

		// последовательный первый запуск задач, если они требуют этого
		// запускаются без таймаута, т.к. требуется, чтобы они полностью завершились перед запуском системы
		if !task.Startup() {
			continue
		}

		if err := task.Do(ctx); err != nil {
			return err
		}
	}

	return nil
}
