package mrview

import (
	"context"
)

type (
	// Validator - предоставляет методы для регистрации пользовательских тегов валидации
	// и проверки структур на соответствие заданным правилам.
	// Используется для валидации входных данных перед их обработкой.
	Validator interface {
		Register(tagName string, fn func(value string) bool) error
		Validate(ctx context.Context, structure any) error
	}
)
