package mrinit

import (
	"sync"

	"github.com/mondegor/go-sysmess/mrerr"
)

type (
	// ErrorManager - comment struct.
	ErrorManager struct {
		generatorID func() string
		caller      func() mrerr.StackTracer
		mu          sync.Mutex
		list        []ManagedError
	}
)

// NewErrorManager - создаёт объект ErrorManager.
func NewErrorManager(generatorID func() string, caller func() mrerr.StackTracer) *ErrorManager {
	return &ErrorManager{
		generatorID: generatorID,
		caller:      caller,
		mu:          sync.Mutex{},
	}
}

// Register - comment method.
func (e *ErrorManager) Register(item ManagedError) {
	e.registerList([]ManagedError{item})
}

// RegisterList - comment method.
func (e *ErrorManager) RegisterList(items []ManagedError) {
	e.registerList(items)
}

func (e *ErrorManager) registerList(items []ManagedError) {
	e.mu.Lock()
	defer e.mu.Unlock()

	for i := range items {
		e.apply(items[i])
		e.list = append(e.list, items[i])
	}
}

func (e *ErrorManager) apply(item ManagedError) {
	generatorID := e.generatorID
	caller := e.caller

	if !item.WithIDGenerator {
		generatorID = nil
	}

	if !item.WithCaller {
		caller = nil
	}

	// WARNING: происходит изменение объекта ошибки
	*item.Err = mrerr.WithExtra(*item.Err, generatorID, caller)
}
