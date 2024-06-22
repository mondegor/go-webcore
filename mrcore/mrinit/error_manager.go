package mrinit

import (
	"sync"

	"github.com/mondegor/go-sysmess/mrerr"
)

type (
	// ErrorManager - comment struct.
	ErrorManager struct {
		extra mrerr.ProtoExtra
		mu    sync.Mutex
		list  []ManagedError
	}
)

// NewErrorManager - создаёт объект ErrorManager.
func NewErrorManager(extra mrerr.ProtoExtra) *ErrorManager {
	return &ErrorManager{
		extra: extra,
		mu:    sync.Mutex{},
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
	extra := e.extra

	if !item.WithCaller {
		extra.Caller = nil
	}

	if !item.WithOnCreated {
		extra.OnCreated = nil
	}

	// WARNING: происходит изменение объекта ошибки
	*item.Err = mrerr.WithExtra(*item.Err, extra)
}
