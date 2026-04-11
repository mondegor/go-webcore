package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"runtime/debug"

	"github.com/mondegor/go-sysmess/errors"
	"github.com/mondegor/go-sysmess/mrlog"
)

// RecoverHandler - middleware для перехвата и обработки panic в HTTP-обработчиках.
//
// Логика работы:
//  1. Оборачивает вызов next.ServeHTTP в recover-блок;
//  2. При панике:
//     - http.ErrAbortHandler не перехватывается (поведение стандартного http.Server);
//     - В debug-режиме stackTrace выводится в stderr;
//     - В production-режиме паника логируется через logger с полным стеком;
//     - Вызывается fatalFunc для отправки клиенту ответа об ошибке.
//
// Важно:
//   - В production-режиме клиент не получает деталей о панике;
//   - fatalFunc может быть nil, тогда отправляется только 500 статус.
func RecoverHandler(logger mrlog.Logger, isDebug bool, fatalFunc http.HandlerFunc) func(next http.Handler) http.Handler {
	if fatalFunc == nil {
		fatalFunc = func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func(ctx context.Context) {
				if rvr := recover(); rvr != nil {
					if rvr == http.ErrAbortHandler { //nolint:errorlint
						// we don't recover http.ErrAbortHandler so the response
						// to the client is aborted, this should not be logged
						panic(rvr)
					}

					errorMessage := fmt.Sprintf("proto='%s', method='%s', url='%s'", r.Proto, r.Method, r.URL)

					if isDebug {
						os.Stderr.Write([]byte(errorMessage))
						os.Stderr.Write([]byte(fmt.Sprintf("; panic: %v\n", rvr)))
						os.Stderr.Write(debug.Stack())
					} else {
						logger.Error(
							ctx,
							"RecoverHandler",
							"error",
							errors.ErrInternalCaughtPanic.New(
								"source", errorMessage,
								"recover", rvr,
								"stack_trace", string(debug.Stack()),
							),
						)
					}

					fatalFunc(w, r)
				}
			}(r.Context())

			next.ServeHTTP(w, r)
		})
	}
}
