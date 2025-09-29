package mrinit

import (
	"context"

	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-sysmess/mrerr/handler"
	"github.com/mondegor/go-sysmess/mrerr/mr"
	"github.com/mondegor/go-sysmess/mrlog"
)

// InitErrorHandler - создаёт объект handler.ErrorHandler.
func InitErrorHandler(logger mrlog.Logger) *handler.ErrorHandler {
	return handler.NewErrorHandler(
		func(ctx context.Context, analyzedKind mrerr.ErrorKind, err error) {
			if analyzedKind == mrerr.ErrorKindUser {
				// 1. пользовательские ошибки: InstantError, ProtoError kind=User
				logger.Debug(ctx, "ErrorHandler", "error", err)

				return
			}

			// 2. пользовательские ошибки с вложенной ошибкой: InstantError kind=User + wrapped err;
			// 3. ошибки: InstantError kind=Internal/System;
			// 4. ProtoError kind=Internal/System (требуется найти место их создания и добавить для них вызов одного из методов New/Wrap);
			// 5. остальные ошибки: которые не были обёрнуты в InstantError (требуется найти место их создания и вложить их в InstantError);
			logger.Error(ctx, "ErrorHandler", "error", err)
		},
		func(err error) error {
			return mr.ErrUnexpectedInternal.Wrap(err)
		},
	)
}
