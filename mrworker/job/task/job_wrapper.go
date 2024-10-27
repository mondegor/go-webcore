package task

import (
	"context"
	"time"

	"github.com/mondegor/go-webcore/mrworker"
)

const (
	defaultCaption = "Task"
	defaultStartup = false
	defaultPeriod  = 30 * time.Second
	defaultTimeout = 60 * time.Second
)

// JobWrapper - обёртка реализующая интерфейс mrworker.Task, используемая
// в планировщике задач, позволяющая вкладывать в себя конкретные работы.
type JobWrapper struct {
	caption string
	startup bool
	period  time.Duration
	timeout time.Duration
	job     mrworker.Job
}

// NewJobWrapper - создаёт объект JobWrapper.
func NewJobWrapper(job mrworker.Job, opts ...Option) *JobWrapper {
	t := &JobWrapper{
		caption: defaultCaption,
		startup: defaultStartup,
		period:  defaultPeriod,
		timeout: defaultTimeout,
		job:     job,
	}

	for _, opt := range opts {
		opt(t)
	}

	return t
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

// Do - исполняет задачу.
func (j *JobWrapper) Do(ctx context.Context) error {
	return j.job.Do(ctx)
}
