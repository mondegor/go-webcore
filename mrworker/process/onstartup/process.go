package onstartup

import (
	"context"
	"sync"
	"time"

	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-sysmess/mrtrace"

	"github.com/mondegor/go-webcore/mrworker"
)

const (
	defaultCaption      = "OnStartup"
	defaultReadyTimeout = 30 * time.Second
)

type (
	// Process - сервис выполнения работы в момент старта приложения. Его полезно использовать,
	// когда работу нужно выполнить после гарантированного запуска остальных процессов.
	Process struct {
		caption      string
		readyTimeout time.Duration
		job          mrworker.Job
		logger       mrlog.Logger
		traceManager mrtrace.ContextManager
		wg           sync.WaitGroup
		done         chan struct{}
	}
)

// NewProcess - создаёт объект Process.
func NewProcess(
	job mrworker.Job,
	logger mrlog.Logger,
	traceManager mrtrace.ContextManager,
	opts ...Option,
) *Process {
	o := options{
		process: &Process{
			caption:      defaultCaption,
			readyTimeout: defaultReadyTimeout,
			job:          job,
			logger:       logger,
			traceManager: traceManager,
			wg:           sync.WaitGroup{},
			done:         make(chan struct{}),
		},
	}

	for _, opt := range opts {
		opt(&o)
	}

	return o.process
}

// Caption - возвращает название сервиса.
func (p *Process) Caption() string {
	return p.caption
}

// ReadyTimeout - возвращает максимальное время, за которое должен быть запущен сервис.
func (p *Process) ReadyTimeout() time.Duration {
	return p.readyTimeout
}

// Start - запуск сервиса выполнения работы при старте приложения.
// Отмена внешнего контекста приведёт к аварийному завершению процесса,
// для корректной остановки следует использовать Shutdown.
// Повторный запуск метода одно и того же объекта не предусмотрен, даже после вызова Shutdown.
func (p *Process) Start(ctx context.Context, ready func()) error {
	p.wg.Add(1)
	defer p.wg.Done()

	p.logger.Debug(ctx, "Starting the startup process...")
	defer p.logger.Debug(ctx, "The startup process has been stopped")

	if err := p.execJob(ctx); err != nil {
		return err
	}

	p.logger.Debug(ctx, "The job of the process is completed")

	if ready != nil {
		ready()
	}

	select {
	case <-p.done:
	case <-ctx.Done():
		p.logger.Debug(ctx, "The startup process detected context 'Done'", "error", ctx.Err())
	}

	return nil
}

// Shutdown - корректная остановка сервиса выполнения работы при старте приложения.
// При повторном вызове метода произойдёт panic.
func (p *Process) Shutdown(ctx context.Context) error {
	p.logger.Debug(ctx, "Shutting down the startup process...")
	close(p.done)

	p.wg.Wait()
	p.logger.Debug(ctx, "The startup process has been shut down")

	return nil
}

func (p *Process) execJob(ctx context.Context) error {
	ctx = p.traceManager.WithGeneratedProcessID(ctx, mrtrace.KeyTaskID)
	p.logger.Debug(ctx, "Execute the job", "job_name", p.Caption())

	if err := p.job.Do(ctx); err != nil {
		return err
	}

	p.logger.Debug(ctx, "The job is completed", "job_name", p.Caption())

	return nil
}
