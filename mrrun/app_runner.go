package mrrun

import (
	"context"

	"github.com/mondegor/go-sysmess/mrlog"
)

type (
	// AppRunner - компонент запуска группы процессов.
	AppRunner struct {
		runner       ProcessRunner
		logger       mrlog.Logger
		traceManager traceManager
	}

	traceManager interface {
		NewContextWithIDs(originalCtx context.Context) context.Context
		WithGeneratedProcessID(ctx context.Context) context.Context
	}
)

// NewAppRunner - создаёт объект AppRunner.
func NewAppRunner(runner ProcessRunner, logger mrlog.Logger, traceManager traceManager) *AppRunner {
	return &AppRunner{
		runner:       runner,
		logger:       logger,
		traceManager: traceManager,
	}
}

// Add - добавляет функции запуска и остановки произвольного процесса.
// Запуск будет осуществлён параллельно с другими добавленными процессами.
func (r *AppRunner) Add(execute func() error, interrupt func(error)) {
	r.runner.Add(execute, interrupt)
}

// AddProcess - добавляет функции запуска и остановки процесса.
// Запуск будет осуществлён параллельно с другими добавленными процессами.
func (r *AppRunner) AddProcess(ctx context.Context, process Process) {
	ex := r.makeExecuter(ctx, process)
	r.runner.Add(ex.Execute, ex.Interrupt)
}

// AddFirstProcess - добавляет функции запуска и остановки процесса.
// Также возвращает канал, по которому будет передано событие, что процесс запущен.
// Запуск будет осуществлён параллельно с другими добавленными процессами.
func (r *AppRunner) AddFirstProcess(ctx context.Context, process Process) (first ProcessSync) {
	ex := r.makeNextExecuter(ctx, process, ProcessSync{})
	r.runner.Add(ex.Execute, ex.Interrupt)

	return ex.Synchronizer
}

// AddNextProcess - добавляет функции запуска и остановки процесса.
// Также возвращает канал, по которому будет передано событие, что процесс запущен.
// Запуск процесса будет осуществлён только при получении события по каналу chPrev (если канал указан).
func (r *AppRunner) AddNextProcess(ctx context.Context, process Process, prev ProcessSync) (next ProcessSync) {
	ex := r.makeNextExecuter(ctx, process, prev)
	r.runner.Add(ex.Execute, ex.Interrupt)

	return ex.Synchronizer
}

// Run - запуск всех добавленных процессов.
// Возвращает ошибку, если хотя бы работа одного процесса прервалась.
func (r *AppRunner) Run() error {
	return r.runner.Run()
}
