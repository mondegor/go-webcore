package mrcore

import (
	"context"

	"github.com/mondegor/go-sysmess/mrerr"
)

type (
	// ErrorHandler - обработчик ошибок.
	ErrorHandler interface {
		Process(ctx context.Context, err error)
	}

	// UsecaseErrorWrapper - помощник для оборачивания UseCase ошибок.
	UsecaseErrorWrapper interface {
		IsNotFoundError(err error) bool
		WrapErrorFailed(err error, source string) error
		WrapErrorNotFoundOrFailed(err error, source string) error
		WrapErrorEntityFailed(err error, entityName string, entityData any) error
		WrapErrorEntityNotFoundOrFailed(err error, entityName string, entityData any) error
	}
)

// errUnexpectedInternal - особая ошибка, в которую заворачивается системой необработанная ошибка.
var errUnexpectedInternal = mrerr.NewProto(
	mrerr.ErrorCodeUnexpectedInternal, mrerr.ErrorKindInternal, "unexpected internal error")

// IsUnexpectedError - проверяется, является ли указанная ошибка необработанной.
// Такая ошибка может появиться только в результате вызова метода CastToAppError.
func IsUnexpectedError(err error) bool {
	return errUnexpectedInternal.Is(err)
}

// PrepareError - приводит указанную ошибку к AppError или если это
// невозможно, то оборачивает её в Internal ошибку, затем возвращает результат.
// (при этом происходит генерация ID и стека вызовов для ошибки, если это необходимо).
// Вызываться должно как можно ближе к тому месту, где произошла непосредственно ошибка.
func PrepareError(err error) error {
	if _, ok := err.(*mrerr.AppError); ok { //nolint:errorlint
		return err
	}

	if proto, ok := err.(*mrerr.ProtoAppError); ok { //nolint:errorlint
		return proto.New()
	}

	return ErrInternal.Wrap(err)
}

// CastToAppError - приводит указанную ошибку к AppError или если это
// невозможно, то оборачивает её в UnexpectedInternal ошибку, затем возвращает результат.
// (при этом исключается генерация ID и стека вызовов для ошибки).
// Используется на поздних этапах обработки ошибок.
func CastToAppError(err error) *mrerr.AppError {
	if appErr, ok := err.(*mrerr.AppError); ok { //nolint:errorlint
		return appErr
	}

	if proto, ok := err.(*mrerr.ProtoAppError); ok { //nolint:errorlint
		return mrerr.Cast(proto)
	}

	// важно, что здесь не происходит генерации ID и стека вызовов
	return errUnexpectedInternal.Wrap(err)
}
