package mrschedule

import (
	"context"
	"time"

	"github.com/mondegor/go-webcore/mrlog"
)

// TaskShell - оболочка реализующая интерфейс mrworker.Task, используемая
// в планировщике задач, позволяющая вкладывать в себя конкретные задачи.
type TaskShell struct {
	caption string
	startup bool
	period  time.Duration
	timeout time.Duration
	doFunc  func(ctx context.Context) error
}

// NewTaskShell - создаёт объект TaskShell.
func NewTaskShell(caption string, startup bool, period, timeout time.Duration, doFunc func(ctx context.Context) error) *TaskShell {
	return &TaskShell{
		caption: caption,
		startup: startup,
		period:  period,
		timeout: timeout,
		doFunc:  doFunc,
	}
}

// Caption - возвращает название задачи.
func (t *TaskShell) Caption() string {
	return t.caption
}

// Startup - необходимо ли стартовать задачу сразу же при инициализации планировщика.
func (t *TaskShell) Startup() bool {
	return t.startup
}

// Period - возвращает периодичность запуска задачи.
func (t *TaskShell) Period() time.Duration {
	return t.period
}

// Timeout - возвращает таймаут выполнения задачи.
func (t *TaskShell) Timeout() time.Duration {
	return t.timeout
}

// Do - исполняет задачу.
func (t *TaskShell) Do(ctx context.Context) error {
	ctx = mrlog.WithContext(ctx, mrlog.Ctx(ctx).With().Str("task", t.caption).Logger())

	return t.doFunc(ctx)
}
