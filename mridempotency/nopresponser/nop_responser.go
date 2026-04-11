package nopresponser

import "net/http"

type (
	// Responser - заглушка (no-op) реализующая интерфейс Responser идемпотентности.
	// Всегда возвращает пустое содержимое и HTTP-код 200 OK.
	// Используется вместе с nopprovider.Provider для тестирования и отладки.
	Responser struct{}
)

// New - создаёт no-op ответчик для использования с nopprovider.Provider.
func New() *Responser {
	return &Responser{}
}

// StatusCode - всегда возвращает HTTP-код ответа 200 (OK).
func (r Responser) StatusCode() int {
	return http.StatusOK
}

// Content - всегда возвращает nil (пустое тело ответа).
func (r Responser) Content() []byte {
	return nil
}
