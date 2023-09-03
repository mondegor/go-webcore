package mrcore

type (
    EventBox interface {
        Emit(message string, args ...any)
    }
)
