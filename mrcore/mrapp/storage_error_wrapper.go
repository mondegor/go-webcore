package mrapp

import (
	"github.com/mondegor/go-sysmess/mrerr"

	"github.com/mondegor/go-webcore/mrcore"
)

type (
	// StorageErrorWrapper - помощник оборачивания перехваченных ошибок
	// в часто используемые ошибки инфраструктурного слоя приложения.
	StorageErrorWrapper struct{}
)

// NewStorageErrorWrapper - создаёт объект StorageErrorWrapper.
func NewStorageErrorWrapper() *StorageErrorWrapper {
	return &StorageErrorWrapper{}
}

// WrapError - возвращает ошибку с указанием источника.
// Если ошибка не mrerr.AppError или не mrerr.ProtoAppError, то она оборачивается в mrcore.ErrStorageQueryFailed.
func (h *StorageErrorWrapper) WrapError(err error, source string) error {
	return h.wrapError(err, "source", source)
}

// WrapErrorEntity - возвращает ошибку с указанием сущности и её данных.
// Если ошибка не mrerr.AppError или не mrerr.ProtoAppError, то она оборачивается в mrcore.ErrStorageQueryFailed.
func (h *StorageErrorWrapper) WrapErrorEntity(err error, entityName string, entityData any) error {
	return h.wrapError(err, entityName, entityData)
}

func (h *StorageErrorWrapper) wrapError(err error, name string, data any) error {
	return mrcore.CastToAppError(
		err,
		func(err error) *mrerr.AppError {
			return mrcore.ErrStorageQueryFailed.Wrap(err)
		},
	).WithAttr(name, data)
}
