package mrrun

import (
	"context"
)

type (
	// AppRunner - компонент запуска группы процессов.
	AppRunner struct {
		runner      ProcessRunner
		chProcesses []chan struct{}
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
// А также возвращает chFirst, при помощи которого можно гарантировать, что отслеживание сигналов системы запущено первым процессом.
// Это необходимо для корректной (graceful) остановки приложения.
func (r *AppRunner) AddSignalHandler(ctx context.Context) (ctxWithCancel context.Context, chFirst chan struct{}) {
	ctx, executer := MakeSignalHandler(ctx)
	r.runner.Add(executer.Execute, executer.Interrupt)
	r.chProcesses = append(r.chProcesses, executer.StartedOk)

	return ctx, executer.StartedOk
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
func (r *AppRunner) AddFirstProcess(ctx context.Context, process Process) (chFirst chan struct{}) {
	executer := MakeNextExecuter(ctx, process, nil)
	r.runner.Add(executer.Execute, executer.Interrupt)
	r.chProcesses = append(r.chProcesses, executer.StartedOk)

	return executer.StartedOk
}

// AddNextProcess - добавляет функции запуска и остановки процесса.
// Также возвращает канал, по которому будет передано событие, что процесс запущен.
// Запуск процесса будет осуществлён только при получении события по каналу chPrev (если канал указан).
func (r *AppRunner) AddNextProcess(ctx context.Context, process Process, chPrev chan struct{}) (chNext chan struct{}) {
	executer := MakeNextExecuter(ctx, process, chPrev)
	r.runner.Add(executer.Execute, executer.Interrupt)
	r.chProcesses = append(r.chProcesses, executer.StartedOk)

	return executer.StartedOk
}

// Run - запуск всех добавленных процессов.
// Если функция onStartup указана, то её вызов произойдёт только после того, как
// все процессы будут запущены (но проверяться будут только те процессы, которым созданы каналы запуска).
// Возвращает ошибку, если хотя бы работа одного процесса прервалась.
func (r *AppRunner) Run(onStartup ...func()) error {
	if len(onStartup) > 0 {
		r.waitingForStartup(onStartup[0])
	}

	return r.runner.Run()
}

// waitingForStartup - ожидание запуска всех процессов для перевода приложения в запущенное состояние.
func (r *AppRunner) waitingForStartup(onStartup func()) {
	if len(r.chProcesses) == 0 {
		onStartup()

		return
	}

	go func() {
		for _, chProcess := range r.chProcesses {
			if chProcess != nil {
				<-chProcess
			}
		}

		onStartup()
	}()
}
