package schedule

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
	// defaultCaption - название планировщика по умолчанию.
	defaultCaption = "TaskScheduler"

	// defaultReadyTimeout - таймаут готовности планировщика по умолчанию.
	defaultReadyTimeout = 30 * time.Second
)

type (
	// TaskScheduler - многопоточный сервис запуска задач по расписанию.
	//
	// Принцип работы:
	//  1. При старте последовательно запускает задачи с task.Startup() == true;
	//  2. Для каждой задачи создаёт отдельный воркер с таймером (task.Period());
	//  3. Воркеры выполняют task.Do(ctx) периодически или по сигналу task.SignalDo();
	//  4. При ошибках вызывает errorHandler, не останавливая другие задачи;
	//
	// Каждая задача выполняется в отдельной горутине с своим таймером.
	TaskScheduler struct {
		caption      string
		readyTimeout time.Duration
		tasks        []mrworker.Task
		errorHandler errors.Handler
		logger       mrlog.Logger
		traceManager mrtrace.ContextManager
		wg           sync.WaitGroup
		done         chan struct{}
	}
)

var (
	// ErrInternalNoTasks - ошибка при отсутствии задач для запуска.
	ErrInternalNoTasks = errors.NewInternalProto("no tasks to start for the task scheduler")

	// ErrInternalZeroParam - ошибка при нулевых period и timeout задачи (attrs: param_name, task_name).
	ErrInternalZeroParam = errors.NewInternalProto("task has zero param for the task scheduler")
)

// NewTaskScheduler - создаёт планировщик задач.
func NewTaskScheduler(
	errorHandler errors.Handler,
	logger mrlog.Logger,
	traceManager mrtrace.ContextManager,
	opts ...Option,
) *TaskScheduler {
	o := options{
		scheduler: &TaskScheduler{
			caption:      defaultCaption,
			readyTimeout: defaultReadyTimeout,

			errorHandler: errorHandler,
			logger:       logger,
			traceManager: traceManager,

			wg:   sync.WaitGroup{},
			done: make(chan struct{}),
		},
	}

	for _, opt := range opts {
		opt(&o)
	}

	if o.captionPrefix != "" {
		o.scheduler.caption = o.captionPrefix + o.scheduler.caption
	}

	return o.scheduler
}

// Caption - возвращает название планировщика задач в свободной форме.
func (p *TaskScheduler) Caption() string {
	return p.caption
}

// ReadyTimeout - возвращает максимальное время, за которое должен быть запущен планировщик со всеми его задачами.
func (p *TaskScheduler) ReadyTimeout() time.Duration {
	return p.readyTimeout
}

// Start - запуск планировщика задач.
//
// Процесс работы:
//  1. Проверяет наличие задач (ошибка ErrInternalNoTasks если пусто);
//  2. Последовательно запускает задачи с task.Startup() == true;
//  3. Для каждой задачи создаёт воркер с таймером (task.Period());
//  4. Воркеры выполняют task.Do(ctx) с таймаутом (task.Timeout());
//  5. Задачи выполняются периодически или по сигналу task.SignalDo();
//  6. При ошибках вызывает errorHandler, не останавливая другие задачи;
//
// Важно:
//   - Отмена внешнего контекста приведёт к завершению всех воркеров;
//   - Для корректной остановки используйте Shutdown;
//   - Повторный запуск того же объекта не поддерживается.
func (p *TaskScheduler) Start(ctx context.Context, ready func()) error {
	p.wg.Add(1)
	defer p.wg.Done()

	p.logger.Debug(ctx, "Starting the task scheduler...")
	defer p.logger.Debug(ctx, "The task scheduler has been stopped")

	if err := p.startup(ctx); err != nil {
		return err
	}

	wgWorkers := sync.WaitGroup{}

	for i := range p.tasks {
		wgWorkers.Add(1)

		go func(ctx context.Context, task mrworker.Task) {
			defer wgWorkers.Done()

			ctx = p.traceManager.WithGeneratedProcessID(ctx, mrtrace.KeyWorkerID)
			p.logger.Info(ctx, "Starting the worker", "task_name", task.Caption())

			ticker := time.NewTicker(task.Period())

			defer func() {
				ticker.Stop()

				if rvr := recover(); rvr != nil {
					p.errorHandler.Handle(
						ctx,
						errors.ErrInternalCaughtPanic.New(
							"source", "task worker: "+task.Caption(),
							"recover", rvr,
							"stack_trace", string(debug.Stack()),
						),
					)
				}
			}()

			for {
				select {
				case <-p.done:
					return
				case <-ctx.Done():
					p.logger.Debug(ctx, "The task worker detected context 'Done'", "task_name", task.Caption(), "error", ctx.Err())

					return
				case <-task.SignalDo():
					p.logger.Debug(ctx, "signalDo event", "task_name", task.Caption())
				case <-ticker.C:
					p.logger.Debug(ctx, "ticker.C event", "task_name", task.Caption())
				}

				ticker.Reset(task.Period())

				if err := p.execTask(ctx, task); err != nil {
					p.errorHandler.Handle(ctx, err)
				}
			}
		}(ctx, p.tasks[i])
	}

	if ready != nil {
		ready()
	}

	wgWorkers.Wait()

	return nil
}

// Shutdown - корректная остановка планировщика задач.
// Останавливает все воркеры и ожидает их завершения.
//
// Важно: при повторном вызове произойдёт panic (закрытие закрытого канала done).
func (p *TaskScheduler) Shutdown(ctx context.Context) error {
	p.logger.Debug(ctx, "Shutting down the task scheduler...")
	close(p.done)

	p.wg.Wait()
	p.logger.Debug(ctx, "The task scheduler has been shut down")

	return nil
}

func (p *TaskScheduler) startup(ctx context.Context) error {
	if len(p.tasks) == 0 {
		return ErrInternalNoTasks.New("scheduler_name", p.caption)
	}

	// запуск задач на этапе старта планировщика выполняется последовательно
	for _, task := range p.tasks {
		if task.Period() == 0 {
			return ErrInternalZeroParam.New("param_name", "period", "task_name", task.Caption())
		}

		if task.Timeout() == 0 {
			return ErrInternalZeroParam.New("param_name", "timeout", "task_name", task.Caption())
		}

		if !task.Startup() {
			continue
		}

		if err := p.execTask(ctx, task); err != nil {
			return err
		}

		if ctx.Err() != nil {
			p.logger.Debug(ctx, "Task scheduler detected context 'Done'", "error", ctx.Err())

			return ctx.Err()
		}
	}

	return nil
}

func (p *TaskScheduler) execTask(ctx context.Context, task mrworker.Task) error {
	ctx = p.traceManager.WithGeneratedProcessID(ctx, mrtrace.KeyTaskID)
	p.logger.Info(ctx, "Execute the task", "task_name", task.Caption())

	ctx, cancel := context.WithTimeout(ctx, task.Timeout())
	defer cancel()

	if err := task.Do(ctx); err != nil {
		return err
	}

	p.logger.Debug(ctx, "The task is completed", "task_name", task.Caption())

	return nil
}
