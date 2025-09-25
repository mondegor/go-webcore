package onstartup

import (
	"context"
	"sync"
	"time"

	"github.com/mondegor/go-sysmess/mrlog"

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
		wgMain       sync.WaitGroup
		done         chan struct{}
	}
)

// NewProcess - создаёт объект Process.
func NewProcess(job mrworker.Job, logger mrlog.Logger, opts ...Option) *Process {
	p := &Process{
		caption:      defaultCaption,
		readyTimeout: defaultReadyTimeout,
		job:          job,
		logger:       logger,
		wgMain:       sync.WaitGroup{},
		done:         make(chan struct{}),
	}

	for _, opt := range opts {
		opt(p)
	}

	return p
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
// Повторный запуск метода одно и того же объекта не предусмотрен, даже после вызова Shutdown.
func (p *Process) Start(ctx context.Context, ready func()) error {
	p.wgMain.Add(1)
	defer p.wgMain.Done()

	p.logger.Debug(ctx, "Starting the process...")
	defer p.logger.Debug(ctx, "The process has been stopped")

	if err := p.job.Do(ctx); err != nil {
		return err
	}

	p.logger.Debug(ctx, "The job of the process is completed")

	if ready != nil {
		ready()
	}

	<-p.done

	return nil
}

// Shutdown - корректная остановка сервиса выполнения работы при старте приложения.
func (p *Process) Shutdown(ctx context.Context) error {
	p.logger.Debug(ctx, "Shutting down the process...")
	close(p.done)

	p.wgMain.Wait()
	p.logger.Debug(ctx, "The process has been shut down")

	return nil
}
