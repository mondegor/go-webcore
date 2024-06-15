package mrcoreerr

import (
	"errors"

	"github.com/mondegor/go-webcore/mrcore"
)

type (
	// UsecaseErrorWrapper - comment struct.
	UsecaseErrorWrapper struct{}
)

// Make sure the Image conforms with the mrcore.UsecaseErrorWrapper interface.
var _ mrcore.UsecaseErrorWrapper = (*UsecaseErrorWrapper)(nil)

// NewUsecaseErrorWrapper - создаёт объект UsecaseErrorWrapper.
func NewUsecaseErrorWrapper() *UsecaseErrorWrapper {
	return &UsecaseErrorWrapper{}
}

// IsNotFoundError - comment method.
func (h *UsecaseErrorWrapper) IsNotFoundError(err error) bool {
	return errors.Is(err, mrcore.ErrStorageNoRowFound) ||
		errors.Is(err, mrcore.ErrStorageRowsNotAffected)
}

// WrapErrorFailed - comment method.
func (h *UsecaseErrorWrapper) WrapErrorFailed(err error, source string) error {
	return h.wrapErrorFailed(err, "source", source)
}

// WrapErrorNotFoundOrFailed - comment method.
func (h *UsecaseErrorWrapper) WrapErrorNotFoundOrFailed(err error, source string) error {
	if h.IsNotFoundError(err) {
		return mrcore.ErrUseCaseEntityNotFound.Wrap(err)
	}

	return h.wrapErrorFailed(err, "source", source)
}

// WrapErrorEntityFailed - comment method.
func (h *UsecaseErrorWrapper) WrapErrorEntityFailed(err error, entityName string, entityData any) error {
	return h.wrapErrorFailed(err, entityName, entityData)
}

// WrapErrorEntityNotFoundOrFailed - comment method.
func (h *UsecaseErrorWrapper) WrapErrorEntityNotFoundOrFailed(err error, entityName string, entityData any) error {
	if h.IsNotFoundError(err) {
		return mrcore.ErrUseCaseEntityNotFound.Wrap(err)
	}

	return h.wrapErrorFailed(err, entityName, entityData)
}

func (h *UsecaseErrorWrapper) wrapErrorFailed(err error, name string, data any) error {
	wrapper := mrcore.ErrUseCaseTemporarilyUnavailable

	if errors.Is(err, mrcore.ErrStorageQueryFailed) {
		wrapper = mrcore.ErrUseCaseOperationFailed
	}

	return wrapper.Wrap(err).WithAttr(name, data)
}
