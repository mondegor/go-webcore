package mrrun

import (
	"sync/atomic"
)

type (
	// AppHealth - компонент управления работоспособностью приложения.
	// Отслеживает состояние запуска (started) и готовности (ready) приложения.
	// Используется health-пробами для reporting состояния приложения.
	//
	// Потокобезопасен благодаря использованию atomic-операций.
	AppHealth struct {
		// isStarted - флаг, указывающий что приложение было запущено.
		isStarted atomic.Bool

		// isReady - флаг, указывающий что приложение готово принимать запросы.
		isReady atomic.Bool
	}
)

// NewAppHealth - создаёт компонент отслеживания работоспособности приложения.
func NewAppHealth() *AppHealth {
	return &AppHealth{}
}

// StartupCompleted - переводит приложение в состояние "запущено и готово".
// Вызывается после успешного запуска всех основных компонентов приложения.
func (a *AppHealth) StartupCompleted() {
	a.isStarted.Store(true)
	a.isReady.Store(true)
}

// IsStarted - сообщает, было ли запущено приложение.
func (a *AppHealth) IsStarted() bool {
	return a.isStarted.Load()
}

// IsReady - сообщает, было ли запущено приложение и готово ли оно принимать запросы.
func (a *AppHealth) IsReady() bool {
	return a.isStarted.Load() && a.isReady.Load()
}

// Pause - временно приостанавливает приём запросов приложением.
func (a *AppHealth) Pause() {
	a.isReady.Store(false)
}

// Continue - возобновляет приём запросов после приостановки.
func (a *AppHealth) Continue() {
	a.isReady.Store(true)
}
