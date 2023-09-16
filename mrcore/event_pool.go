package mrcore

import "context"

type (
    EventPool interface {
        Listen(eventName string, handlerFunc EventHandlerFunc)
        Emit(eventName string, args EventArgs) error
        EmitWithContext(ctx context.Context, eventName string, args EventArgs) error
        List() []string
        Has(eventName string) bool
        Remove(eventNames ...string)
        Wait()
    }

    EventArgs any

    EventHandlerFunc func(ctx context.Context, args EventArgs)
)
