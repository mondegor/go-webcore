package mrcoreerr

import (
	"context"

	"github.com/mondegor/go-webcore/mrlog"
)

type (
	// ErrorHandler - comment struct.
	ErrorHandler struct{}
)

// NewErrorHandler - создаёт объект ErrorHandler.
func NewErrorHandler() *ErrorHandler {
	return &ErrorHandler{}
}

// Process - comment method.
func (h *ErrorHandler) Process(ctx context.Context, err error) {
	mrlog.Ctx(ctx).Error().Err(err).Msg("mrcoreerr.Log")
}
