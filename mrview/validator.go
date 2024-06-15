package mrview

import (
	"context"
)

type (
	// Validator - интерфейс валидации полей структур, с возможностью регистрации
	// необходимых тегов, которые можно заранее привязывать к полям структур.
	Validator interface {
		Register(tagName string, fn func(value string) bool) error
		Validate(ctx context.Context, structure any) error
	}
)
