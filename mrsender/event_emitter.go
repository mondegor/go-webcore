package mrsender

import (
	"context"
)

const (
	// SourceEventSeparator - разделитель между источником и названием события.
	SourceEventSeparator = ":"

	// DefaultSource - название источника по умолчанию.
	DefaultSource = "EmptySource"
)

// :TODO: заменить object any ---> args ...any

type (
	// EventEmitter - отправитель событий.
	EventEmitter interface {
		Emit(ctx context.Context, eventName string, object any)
	}

	// EventReceiver - получатель событий.
	EventReceiver interface {
		Receive(ctx context.Context, eventName, source string, object any)
	}

	// EventReceiveFunc - получатель событий в виде функции.
	EventReceiveFunc func(ctx context.Context, eventName, source string, object any)
)

// Receive - реализация интерфейса EventReceiver в виде функции для получения события.
func (f EventReceiveFunc) Receive(ctx context.Context, eventName, source string, object any) {
	f(ctx, eventName, source, object)
}
