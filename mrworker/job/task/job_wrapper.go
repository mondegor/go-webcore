package task

import (
	"context"
	"time"

	"github.com/mondegor/go-webcore/mrworker"
)

const (
	// defaultCaption - название задачи по умолчанию.
	defaultCaption = "Task"

	// defaultStartup - значение запуска задачи при старте по умолчанию.
	defaultStartup = false

	// defaultTimeout - таймаут выполнения задачи по умолчанию.
	defaultTimeout = 30 * time.Second
)

// defaultPeriod - периодичность запуска задачи по умолчанию.
var defaultPeriod = mrworker.NewStaticPeriod(60 * time.Second) //nolint:gochecknoglobals

// JobWrapper - обёртка, реализующая интерфейс mrworker.Task для использования
// в планировщике задач (TaskScheduler). Позволяет адаптировать любой mrworker.Job
// к требованиям планировщика, добавляя настройки периода, таймаута и сигналов.
type (
	JobWrapper struct {
		caption        string
		startup        bool
		periodStrategy mrworker.PeriodStrategy
		timeout        time.Duration
		signalDo       <-chan struct{} // signalDo - канал для немедленного запуска задачи
		job            mrworker.Job
	}
)

// NewJobWrapper - создаёт обёртку Task для указанной задачи.
func NewJobWrapper(job mrworker.Job, opts ...Option) *JobWrapper {
	o := options{
		job: &JobWrapper{
			caption:        defaultCaption,
			startup:        defaultStartup,
			periodStrategy: defaultPeriod,
			timeout:        defaultTimeout,
			job:            job,
		},
	}

	for _, opt := range opts {
		opt(&o)
	}

	if o.captionPrefix != "" {
		o.job.caption = o.captionPrefix + o.job.caption
	}

	return o.job
}

// Caption - возвращает название задачи.
func (j *JobWrapper) Caption() string {
	return j.caption
}

// Startup - сообщает, если задачу нужно выполнить сразу же при старте планировщика.
func (j *JobWrapper) Startup() bool {
	return j.startup
}

// Period - возвращает периодичность запуска задачи.
func (j *JobWrapper) Period() time.Duration {
	return j.periodStrategy.Period()
}

// Timeout - возвращает максимальное время выполнения задачи.
func (j *JobWrapper) Timeout() time.Duration {
	return j.timeout
}

// SignalDo - возвращает канал для немедленного запуска задачи.
func (j *JobWrapper) SignalDo() <-chan struct{} {
	return j.signalDo
}

// Do - выполняет обёрнутую задачу.
func (j *JobWrapper) Do(ctx context.Context) error {
	return j.job.Do(ctx)
}
