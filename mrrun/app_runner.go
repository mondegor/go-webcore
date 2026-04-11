package mrrun

import (
	"context"

	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-sysmess/mrtrace"
)

type (
	// AppRunner - компонент запуска и управления группой процессов (сервисов).
	// Обеспечивает параллельный запуск процессов, их синхронизацию и корректное завершение.
	//
	// Поддерживает несколько режимов добавления процессов:
	//   - Add - простой процесс без синхронизации;
	//   - AddFirstProcess - первый процесс, от которого зависят остальные;
	//   - AddNextProcess - процесс, зависящий от запуска предыдущего процесса.
	AppRunner struct {
		runner       ProcessRunner
		logger       mrlog.Logger
		traceManager mrtrace.ContextManager
	}
)

// NewAppRunner - создаёт менеджер запуска процессов.
func NewAppRunner(runner ProcessRunner, logger mrlog.Logger, traceManager mrtrace.ContextManager) *AppRunner {
	return &AppRunner{
		runner:       runner,
		logger:       logger,
		traceManager: traceManager,
	}
}

// Add - добавляет процесс через функции запуска и остановки.
// Процесс запускается параллельно с другими добавленными процессами.
func (r *AppRunner) Add(execute func() error, interrupt func(error)) {
	r.runner.Add(execute, interrupt)
}

// AddProcess - добавляет процесс, реализующий интерфейс Process.
// Процесс запускается параллельно с другими, без ожидания запуска других.
func (r *AppRunner) AddProcess(ctx context.Context, process Process) {
	ex := r.makeExecuter(ctx, process)
	r.runner.Add(ex.Execute, ex.Interrupt)
}

// AddFirstProcess - добавляет первый процесс, от которого могут зависеть остальные.
// Возвращает канал синхронизации (ProcessSync), по которому другие процессы
// могут ожидать завершения запуска этого процесса.
//
// Процесс запускается параллельно с другими, но другие процессы могут
// дождаться его готовности через возвращаемый ProcessSync.
func (r *AppRunner) AddFirstProcess(ctx context.Context, process Process) (first ProcessSync) {
	ex := r.makeNextExecuter(ctx, process, ProcessSync{})
	r.runner.Add(ex.Execute, ex.Interrupt)

	return ex.Synchronizer
}

// AddNextProcess - добавляет процесс, зависящий от запуска предыдущего процесса.
// Процесс начнёт запуск только после получения сигнала через канал prev.
// Если prev не указан (пустой ProcessSync), процесс запустится немедленно.
//
// Возвращает канал синхронизации (ProcessSync) для процессов, зависящих от этого.
func (r *AppRunner) AddNextProcess(ctx context.Context, process Process, prev ProcessSync) (next ProcessSync) {
	ex := r.makeNextExecuter(ctx, process, prev)
	r.runner.Add(ex.Execute, ex.Interrupt)

	return ex.Synchronizer
}

// Run - запускает все добавленные процессы параллельно.
// Блокируется до завершения всех процессов или до ошибки любого процесса.
//
// Возвращает ошибку первого процесса, который завершился с ошибкой.
func (r *AppRunner) Run() error {
	return r.runner.Run()
}
