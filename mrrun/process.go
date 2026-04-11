package mrrun

import (
	"context"
	"time"
)

type (
	// Process - интерфейс процесса (сервиса) для параллельного запуска и остановки.
	//
	// Жизненный цикл процесса:
	//  1. Вызывается ReadyTimeout() для получения таймаута готовности;
	//  2. Вызывается Start(ctx, ready) - процесс должен вызвать ready() после запуска;
	//  3. При завершении вызывается Shutdown(ctx) для корректной остановки;
	Process interface {
		Caption() string
		ReadyTimeout() time.Duration
		Start(ctx context.Context, ready func()) error
		Shutdown(ctx context.Context) error
	}

	// ProcessRunner - интерфейс исполнителя процессов.
	// Управляет добавлением и параллельным запуском процессов,
	// а также обеспечивает их корректное прерывание при ошибке.
	ProcessRunner interface {
		Add(execute func() error, interrupt func(error))
		Run() error
	}

	// ProcessSync - канал синхронизации для передачи сигнала готовности процесса.
	// Используется для последовательного запуска процессов, когда следующий
	// процесс ждёт готовности предыдущего.
	ProcessSync struct {
		Caption      string
		readyTimeout time.Duration

		// ready - канал, который закрывается при готовности процесса.
		// Закрытие канала сигнализирует следующим процессам о готовности.
		ready chan struct{}
	}
)
