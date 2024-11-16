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

// AddProcess - добавляет функции запуска и остановки процесса.
// Запуск будет осуществлён параллельно с другими добавленными процессами.
func (r *AppRunner) AddProcess(ctx context.Context, process Process) {
	executer := MakeExecuter(ctx, process)
	r.runner.Add(executer.Execute, executer.Interrupt)
}

// AddFirstProcess - добавляет функции запуска и остановки процесса.
// Также возвращает канал, по которому будет передано событие, что процесс запущен.
// Запуск будет осуществлён параллельно с другими добавленными процессами.
func (r *AppRunner) AddFirstProcess(ctx context.Context, process Process) (first StartingProcess) {
	executer := MakeNextExecuter(ctx, process, StartingProcess{})
	r.runner.Add(executer.Execute, executer.Interrupt)

	return executer.Starting
}

// AddNextProcess - добавляет функции запуска и остановки процесса.
// Также возвращает канал, по которому будет передано событие, что процесс запущен.
// Запуск процесса будет осуществлён только при получении события по каналу chPrev (если канал указан).
func (r *AppRunner) AddNextProcess(ctx context.Context, process Process, prev StartingProcess) (next StartingProcess) {
	executer := MakeNextExecuter(ctx, process, prev)
	r.runner.Add(executer.Execute, executer.Interrupt)

	return executer.Starting
}

// Run - запуск всех добавленных процессов.
// Возвращает ошибку, если хотя бы работа одного процесса прервалась.
func (r *AppRunner) Run() error {
	return r.runner.Run()
}
