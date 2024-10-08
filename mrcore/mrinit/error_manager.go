package mrinit

import (
	"sync"

	"github.com/mondegor/go-sysmess/mrerr"
)

type (
	// ErrorManager - менеджер ошибок для централизованного
	// подключения дополнительных свойств ошибкам ProtoAppError.
	ErrorManager struct {
		extra mrerr.ProtoExtra
		mu    sync.Mutex
		list  []EnrichedError
	}
)

// NewErrorManager - создаёт объект ErrorManager.
func NewErrorManager(extra mrerr.ProtoExtra) *ErrorManager {
	return &ErrorManager{
		extra: extra,
		mu:    sync.Mutex{},
	}
}

// Register - регистрирует указанную ошибку с её дополнительными свойствами.
func (e *ErrorManager) Register(item EnrichedError) {
	e.registerList([]EnrichedError{item})
}

// RegisterList - регистрирует список указанных ошибок с их дополнительными свойствами.
func (e *ErrorManager) RegisterList(items []EnrichedError) {
	e.registerList(items)
}

func (e *ErrorManager) registerList(items []EnrichedError) {
	e.mu.Lock()
	defer e.mu.Unlock()

	for i := range items {
		e.apply(items[i])
		e.list = append(e.list, items[i])
	}
}

func (e *ErrorManager) apply(item EnrichedError) {
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
