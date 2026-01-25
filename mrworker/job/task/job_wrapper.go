package task

import (
	"context"
	"time"

	"github.com/mondegor/go-webcore/mrworker"
)

const (
	defaultCaption = "Task"
	defaultStartup = false
	defaultPeriod  = 60 * time.Second
	defaultTimeout = 30 * time.Second
)

// JobWrapper - обёртка реализующая интерфейс mrworker.Task, используемая
// в планировщике задач, позволяющая вкладывать в себя конкретные работы.
type (
	JobWrapper struct {
		caption  string
		startup  bool
		period   time.Duration
		timeout  time.Duration
		signalDo <-chan struct{}
		job      mrworker.Job
	}
)

// NewJobWrapper - создаёт объект JobWrapper.
func NewJobWrapper(job mrworker.Job, opts ...Option) *JobWrapper {
	o := options{
		job: &JobWrapper{
			caption: defaultCaption,
			startup: defaultStartup,
			period:  defaultPeriod,
			timeout: defaultTimeout,
			job:     job,
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

// Startup - необходимо ли стартовать задачу сразу же при инициализации планировщика.
func (j *JobWrapper) Startup() bool {
	return j.startup
}

// Period - возвращает периодичность запуска задачи.
func (j *JobWrapper) Period() time.Duration {
	return j.period
}

// Timeout - возвращает таймаут выполнения задачи.
func (j *JobWrapper) Timeout() time.Duration {
	return j.timeout
}

// SignalDo - возвращает сигнал, о том, что можно немедленно
// запускать задачу не дожидаясь завершения Period.
func (j *JobWrapper) SignalDo() <-chan struct{} {
	return j.signalDo
}

// Do - исполняет задачу.
func (j *JobWrapper) Do(ctx context.Context) error {
	return j.job.Do(ctx)
}
