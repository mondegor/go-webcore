package onstartup

import (
	"context"
	"time"

	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrworker"
)

const (
	defaultCaption      = "OnStartup"
	defaultReadyTimeout = 60 * time.Second
)

// Process - сервис выполнения работы в момент старта приложения. Его полезно использовать,
// когда работу нужно выполнить после гарантированного запуска остальных процессов.
type Process struct {
	caption      string
	readyTimeout time.Duration
	job          mrworker.Job
	done         chan struct{}
}

// NewProcess - создаёт объект Process.
func NewProcess(job mrworker.Job, opts ...Option) *Process {
	p := &Process{
		caption:      defaultCaption,
		readyTimeout: defaultReadyTimeout,
		job:          job,
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
func (p *Process) Start(ctx context.Context, ready func()) error {
	mrlog.Ctx(ctx).Debug().Msg("Starting the process...")

	if err := p.job.Do(ctx); err != nil {
		return err
	}

	mrlog.Ctx(ctx).Debug().Msg("The job of the process is completed")

	if ready != nil {
		ready()
	}

	<-p.done
	mrlog.Ctx(ctx).Debug().Msg("The process has been stopped")

	return nil
}

// Shutdown - корректная остановка сервиса выполнения работы при старте приложения.
func (p *Process) Shutdown(ctx context.Context) error {
	mrlog.Ctx(ctx).Debug().Msg("Shutting down the process...")
	close(p.done)

	return nil
}
