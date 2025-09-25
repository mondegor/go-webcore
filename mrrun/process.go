package mrrun

import (
	"context"
	"time"
)

type (
	// Process - процесс (сервис) который можно запускать и
	// останавливать параллельно с другими процессами.
	Process interface {
		Caption() string
		ReadyTimeout() time.Duration
		Start(ctx context.Context, ready func()) error
		Shutdown(ctx context.Context) error
	}

	// ProcessRunner - запускатель процессов.
	ProcessRunner interface {
		Add(execute func() error, interrupt func(error))
		Run() error
	}

	// ProcessSync - содержит информацию о процессе находящемся в момент запуска.
	ProcessSync struct {
		Caption      string
		readyTimeout time.Duration
		ready        chan struct{}
	}
)
