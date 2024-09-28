package mrrun

import (
	"sync/atomic"
)

type (
	// AppHealth - компонент управление работоспособностью приложения.
	AppHealth struct {
		isStarted atomic.Bool
		isReady   atomic.Bool
	}
)

// NewAppHealth - создаёт объект AppRunner.
func NewAppHealth() *AppHealth {
	return &AppHealth{}
}

// StartupCompleted - перевод приложения в запущенное состояние.
func (a *AppHealth) StartupCompleted() {
	a.isStarted.Store(true)
	a.isReady.Store(true)
}

// IsStarted - возвращает запущено ли приложение.
func (a *AppHealth) IsStarted() bool {
	return a.isStarted.Load()
}

// IsReady - возвращает готово ли приложение к приёму запросов.
func (a *AppHealth) IsReady() bool {
	return a.isStarted.Load() && a.isReady.Load()
}

// Pause - временное прекращение приложением приёма запросов.
func (a *AppHealth) Pause() {
	a.isReady.Store(false)
}

// Continue - возобновление приложением приёма запросов.
func (a *AppHealth) Continue() {
	a.isReady.Store(true)
}
