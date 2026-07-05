package middleware

import (
	"net/http"

	"github.com/mondegor/go-webcore/mrserver"
)

// HandlerAdapter - адаптирует HttpHandlerFunc к стандартному http.HandlerFunc.
//
// Логика работы:
//  1. Вызывает следующий обработчик next(w, r);
//  2. Если next возвращает ошибку, отправляет её через errSender.SendError;
//  3. Если ошибки нет, ответ уже записан следующим обработчиком.
func HandlerAdapter(errSender mrserver.ErrorResponseSender) func(next mrserver.HttpHandlerFunc) http.HandlerFunc {
	return func(next mrserver.HttpHandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if err := next(w, r); err != nil {
				errSender.SendError(w, r, err)
			}
		}
	}
}
