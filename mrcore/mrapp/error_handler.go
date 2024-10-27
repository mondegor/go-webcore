package mrapp

import (
	"context"

	"github.com/mondegor/go-sysmess/mrerr"

	"github.com/mondegor/go-webcore/mrcore"
)

type (
	// ErrorHandler - обработчик ошибок приложения.
	ErrorHandler struct {
		hook func(ctx context.Context, errType mrcore.AnalyzedErrorType, err *mrerr.AppError)
	}
)

// NewErrorHandler - создаёт объект ErrorHandler.
func NewErrorHandler() *ErrorHandler {
	return &ErrorHandler{}
}

// NewErrorHandlerWithHook - создаёт объект ErrorHandler с указанным хук функцией.
func NewErrorHandlerWithHook(hook func(ctx context.Context, errType mrcore.AnalyzedErrorType, err *mrerr.AppError)) *ErrorHandler {
	return &ErrorHandler{
		hook: hook,
	}
}

// Perform - версия метода PerformWithCommit без вызова коммит функции.
func (h *ErrorHandler) Perform(ctx context.Context, err error) {
	h.PerformWithCommit(ctx, err, nil)
}

// PerformWithCommit - анализирует ошибку, при необходимости оборачивает её
// в ErrUnexpectedInternal и вызывает хук функцию, а затем коммит функцию.
// В результате вызова этих функций ошибка может быть, например, залогирована
// и использована при формировании ответа серверу.
func (h *ErrorHandler) PerformWithCommit(ctx context.Context, err error, commit func(errType mrcore.AnalyzedErrorType, err *mrerr.AppError)) {
	errType := h.analyzeError(err)

	// ошибки с типом AnalyzedErrorTypeUndefined оборачиваются в ErrUnexpectedInternal
	appErr := mrcore.CastToAppError(
		err,
		func(err error) *mrerr.AppError {
			return mrcore.ErrUnexpectedInternal.Wrap(err)
		},
	)

	if h.hook != nil {
		h.hook(ctx, errType, appErr)
	}

	if commit != nil {
		commit(errType, appErr)
	}
}

func (h *ErrorHandler) analyzeError(err error) mrcore.AnalyzedErrorType {
	var isUserError bool

	nestedErr := err

	for {
		if appErr, ok := nestedErr.(*mrerr.AppError); ok { //nolint:errorlint
			if appErr.Kind() == mrerr.ErrorKindSystem {
				return mrcore.AnalyzedErrorTypeSystem
			}

			if appErr.Kind() == mrerr.ErrorKindInternal {
				return mrcore.AnalyzedErrorTypeInternal
			}

			// фиксируется, что изначально ошибка пользовательская
			isUserError = true

			// выбирается причина пользовательской ошибки если такая существует
			if nestedErr = appErr.Unwrap(); nestedErr != nil {
				continue
			}
		}

		break
	}

	// если это пользовательская ошибка, то она точно не содержит internal и system ошибок
	if isUserError {
		return mrcore.AnalyzedErrorTypeUser
	}

	if appErr, ok := err.(*mrerr.ProtoAppError); ok { //nolint:errorlint
		if appErr.Kind() == mrerr.ErrorKindSystem {
			return mrcore.AnalyzedErrorTypeProtoSystem
		}

		if appErr.Kind() == mrerr.ErrorKindInternal {
			return mrcore.AnalyzedErrorTypeProtoInternal
		}

		return mrcore.AnalyzedErrorTypeProtoUser
	}

	return mrcore.AnalyzedErrorTypeUndefined
}
