package mrcore

import "context"

type (
	Validator interface {
		Register(tagName string, fn ValidatorTagNameFunc) error
		Validate(ctx context.Context, structure any) error
	}

	ValidatorTagNameFunc func(value string) bool
)
