package decorator

import (
	"context"
	"strings"

	"github.com/mondegor/go-webcore/mrsender"
)

type (
	// SourceEmitter - обёртка отправителя событий добавляющая первоисточник.
	SourceEmitter struct {
		eventEmitter mrsender.EventEmitter
		source       string
	}
)

// NewSourceEmitter - создаёт объект SourceEmitter.
func NewSourceEmitter(eventEmitter mrsender.EventEmitter, source string) *SourceEmitter {
	if source == "" {
		source = mrsender.DefaultSource
	}

	return &SourceEmitter{
		eventEmitter: eventEmitter,
		source:       source,
	}
}

// Emit - отправляет указанное событие добавляя первоисточник.
func (e *SourceEmitter) Emit(ctx context.Context, eventName string, object any) {
	separator := mrsender.SourceEventSeparator

	if strings.Contains(eventName, separator) {
		separator = "/"
	}

	e.eventEmitter.Emit(ctx, e.source+separator+eventName, object)
}
