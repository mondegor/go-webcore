package mrcoreerr

import (
	"errors"

	"github.com/mondegor/go-webcore/mrcore"
)

type (
	// UseCaseErrorWrapper - помощник для оборачивания ошибок в часто используемые UseCase ошибки.
	UseCaseErrorWrapper struct{}
)

// Make sure the Image conforms with the mrcore.UseCaseErrorWrapper interface.
var _ mrcore.UseCaseErrorWrapper = (*UseCaseErrorWrapper)(nil)

// NewUseCaseErrorWrapper - создаёт объект UseCaseErrorWrapper.
func NewUseCaseErrorWrapper() *UseCaseErrorWrapper {
	return &UseCaseErrorWrapper{}
}

// IsNotFoundError - проверяет, является ли ошибка связанной с тем,
// что запрос валидный, но запись не найдена или её изменение не потребовалось.
func (h *UseCaseErrorWrapper) IsNotFoundError(err error) bool {
	return errors.Is(err, mrcore.ErrStorageNoRowFound) ||
		errors.Is(err, mrcore.ErrStorageRowsNotAffected)
}

// WrapErrorFailed - возвращает ошибку с указанием источника, обёрнутую в
// mrcore.ErrUseCaseTemporarilyUnavailable или mrcore.ErrUseCaseOperationFailed.
func (h *UseCaseErrorWrapper) WrapErrorFailed(err error, source string) error {
	return h.wrapErrorFailed(err, "source", source)
}

// WrapErrorNotFoundOrFailed - возвращает ошибку с указанием источника, обёрнутую в
// mrcore.ErrUseCaseEntityNotFound, mrcore.ErrUseCaseTemporarilyUnavailable или mrcore.ErrUseCaseOperationFailed.
func (h *UseCaseErrorWrapper) WrapErrorNotFoundOrFailed(err error, source string) error {
	if h.IsNotFoundError(err) {
		return mrcore.ErrUseCaseEntityNotFound.Wrap(err)
	}

	return h.wrapErrorFailed(err, "source", source)
}

// WrapErrorEntityFailed - возвращает ошибку с указанием сущности и её данных, обёрнутую в
// mrcore.ErrUseCaseTemporarilyUnavailable или mrcore.ErrUseCaseOperationFailed.
func (h *UseCaseErrorWrapper) WrapErrorEntityFailed(err error, entityName string, entityData any) error {
	return h.wrapErrorFailed(err, entityName, entityData)
}

// WrapErrorEntityNotFoundOrFailed - возвращает ошибку с указанием сущности и её данных, обёрнутую в
// mrcore.ErrUseCaseEntityNotFound, mrcore.ErrUseCaseTemporarilyUnavailable или mrcore.ErrUseCaseOperationFailed.
func (h *UseCaseErrorWrapper) WrapErrorEntityNotFoundOrFailed(err error, entityName string, entityData any) error {
	if h.IsNotFoundError(err) {
		return mrcore.ErrUseCaseEntityNotFound.Wrap(err)
	}

	return h.wrapErrorFailed(err, entityName, entityData)
}

func (h *UseCaseErrorWrapper) wrapErrorFailed(err error, name string, data any) error {
	wrapper := mrcore.ErrUseCaseTemporarilyUnavailable

	if errors.Is(err, mrcore.ErrStorageQueryFailed) {
		wrapper = mrcore.ErrUseCaseOperationFailed
	}

	return wrapper.Wrap(err).WithAttr(name, data)
}
