package mrrun

import (
	"context"
)

type (
	// AppRunner - компонент запуска группы процессов.
	AppRunner struct {
		runner ProcessRunner
	}
)

// NewAppRunner - создаёт объект AppRunner.
func NewAppRunner(runner ProcessRunner) *AppRunner {
	return &AppRunner{
		runner: runner,
	}
}

// Add - добавляет функции запуска и остановки произвольного процесса.
// Запуск будет осуществлён параллельно с другими добавленными процессами.
func (r *AppRunner) Add(execute func() error, interrupt func(error)) {
	r.runner.Add(execute, interrupt)
}

// AddSignalHandler - добавляет функции для отслеживания/прекращения отслеживания
// сигналов системы. Возвращает контекст, в котором установлена его отмена при перехвате системного события.
// Это необходимо для корректной (graceful) остановки приложения.
func (r *AppRunner) AddSignalHandler(ctx context.Context) context.Context {
	ctx, execute, interrupt := MakeSignalHandler(ctx)
	r.runner.Add(execute, interrupt)

	return ctx
}

// AddProcess - добавляет функции запуска и остановки процесса.
// Запуск будет осуществлён параллельно с другими добавленными процессами.
func (r *AppRunner) AddProcess(ctx context.Context, process Process) {
	execute, interrupt := MakeExecuter(ctx, process)
	r.runner.Add(execute, interrupt)
}

// AddHeadProcess - добавляет функции запуска и остановки процесса.
// Также возвращает канал, по которому будет передано событие, что процесс запущен.
// Запуск будет осуществлён параллельно с другими добавленными процессами.
func (r *AppRunner) AddHeadProcess(ctx context.Context, process Process) (chHead chan struct{}) {
	chHead, execute, interrupt := MakeNextExecuter(ctx, process, nil)
	r.runner.Add(execute, interrupt)

	return chHead
}

// AddNextProcess - добавляет функции запуска и остановки процесса.
// Также возвращает канал, по которому будет передано событие, что процесс запущен.
// Запуск процесса будет осуществлён только при получении события по каналу chPrev (если канал указан).
func (r *AppRunner) AddNextProcess(ctx context.Context, process Process, chPrev chan struct{}) (chNext chan struct{}) {
	chNext, execute, interrupt := MakeNextExecuter(ctx, process, chPrev)
	r.runner.Add(execute, interrupt)

	return chNext
}

// Run - запуск всех добавленных ранее процессов.
// Возвращает ошибку, если хотя бы работа одного процесса прервалась.
func (r *AppRunner) Run() error {
	return r.runner.Run()
}
