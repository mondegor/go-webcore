package mrworker

import (
	"context"
	"time"
)

type (
	// Task - интерфейс задачи для планировщика.
	Task interface {
		Caption() string
		Startup() bool
		Period() time.Duration
		Timeout() time.Duration
		Do(ctx context.Context) error
	}
)
