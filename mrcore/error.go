package mrcore

import (
	"context"

	"github.com/mondegor/go-sysmess/mrerr"
)

//go:generate mockgen -source=error.go -destination=./mock/error.go

type (
	// ErrorHandler - обработчик ошибок.
	ErrorHandler interface {
		Handle(ctx context.Context, err error)
		HandleWith(ctx context.Context, err error, extraHandler func(analyzedKind mrerr.ErrorKind, err error))
	}

	// ErrorWrapper - помощник для оборачивания ошибок.
	ErrorWrapper interface {
		WrapError(err error, attrs ...any) error
	}

	// UseCaseErrorWrapper - помощник для оборачивания UseCase ошибок.
	UseCaseErrorWrapper interface {
		IsNotFoundOrNotAffectedError(err error) bool
		WrapErrorFailed(err error, attrs ...any) error
		WrapErrorNotFoundOrFailed(err error, attrs ...any) error
	}
)
