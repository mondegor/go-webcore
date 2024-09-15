package mrrun

import (
	"context"
)

type (
	// Process - процесс (сервис) который можно запускать и
	// останавливать параллельно с другими процессами.
	Process interface {
		Caption() string
		Start(ctx context.Context, ready func()) error
		Shutdown(ctx context.Context) error
	}

	// ProcessRunner - запускатель процессов.
	ProcessRunner interface {
		Add(execute func() error, interrupt func(error))
		Run() error
	}
)
