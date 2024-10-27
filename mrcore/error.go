package mrcore

import (
	"context"

	"github.com/mondegor/go-sysmess/mrerr"
)

const (
	AnalyzedErrorTypeUndefined     AnalyzedErrorType = iota // AnalyzedErrorTypeUndefined - любая ошибка, которая не является ошибкой из списка ниже
	AnalyzedErrorTypeInternal                               // AnalyzedErrorTypeInternal - AppError with kind=ErrorKindInternal
	AnalyzedErrorTypeSystem                                 // AnalyzedErrorTypeSystem - AppError with kind=ErrorKindSystem
	AnalyzedErrorTypeUser                                   // AnalyzedErrorTypeUser - AppError with kind=ErrorKindUser
	AnalyzedErrorTypeProtoInternal                          // AnalyzedErrorTypeProtoInternal - ProtoAppError with kind=ErrorKindInternal
	AnalyzedErrorTypeProtoSystem                            // AnalyzedErrorTypeProtoSystem - ProtoAppError with kind=ErrorKindSystem
	AnalyzedErrorTypeProtoUser                              // AnalyzedErrorTypeProtoUser - ProtoAppError with kind=ErrorKindUser
)

type (
	// AnalyzedErrorType - тип ошибки определённый обработчиком ошибок (ErrorHandler).
	AnalyzedErrorType int8

	// ErrorHandler - обработчик ошибок.
	ErrorHandler interface {
		Perform(ctx context.Context, err error)
		PerformWithCommit(ctx context.Context, err error, commit func(errType AnalyzedErrorType, err *mrerr.AppError))
	}

	// UseCaseErrorWrapper - помощник для оборачивания UseCase ошибок.
	UseCaseErrorWrapper interface {
		IsNotFoundError(err error) bool
		WrapErrorFailed(err error, source string) error
		WrapErrorNotFoundOrFailed(err error, source string) error
		WrapErrorEntityFailed(err error, entityName string, entityData any) error
		WrapErrorEntityNotFoundOrFailed(err error, entityName string, entityData any) error
	}

	// StorageErrorWrapper - помощник для оборачивания Storage ошибок.
	StorageErrorWrapper interface {
		WrapError(err error, source string) error
		WrapErrorEntity(err error, entityName string, entityData any) error
	}
)

// PrepareError - приводит указанную ошибку к AppError или если это
// невозможно, то оборачивает её в ErrInternal ошибку, затем возвращает результат.
// (при этом происходит генерация ID и стека вызовов для ошибки, если в ошибке это предусмотрено).
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

// CastToAppError - приводит указанную ошибку к AppError или если это невозможно,
// то вызывает функцию по умолчанию, указанную в defFunc параметре, затем возвращает результат.
func CastToAppError(err error, defFunc func(err error) *mrerr.AppError) *mrerr.AppError {
	if appErr, ok := err.(*mrerr.AppError); ok { //nolint:errorlint
		return appErr
	}

	if proto, ok := err.(*mrerr.ProtoAppError); ok { //nolint:errorlint
		return mrerr.Cast(proto)
	}

	return defFunc(err)
}
