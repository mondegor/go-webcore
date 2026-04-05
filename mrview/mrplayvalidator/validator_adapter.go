package mrplayvalidator

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/mondegor/go-sysmess/errors"
	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-sysmess/util/xstrings"
)

// go get -u github.com/go-playground/validator/v10

const (
	validatorErrorPrefix           = "Validator_"
	validatorErrorPostfix          = ": {Name}, {Type}, {Value}"
	validatorErrorPostfixWithParam = validatorErrorPostfix + ", {Param}"
	validatorErrorID               = "ValidateError"
)

type (
	// ValidatorAdapter - адаптер валидатора структур и их полей на базе тегов.
	ValidatorAdapter struct {
		validate  *validator.Validate
		logger    mrlog.Logger
		tag2error map[string]errors.UserProtoError
	}
)

var errInternalValidatorTagIsNotFound = errors.NewInternalProto("validator error: tag is empty")

// New - создаёт объект ValidatorAdapter.
func New(logger mrlog.Logger) *ValidatorAdapter {
	validate := validator.New()

	// возвращение в качестве имени поля названия из тега json
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}

		return name
	})

	return &ValidatorAdapter{
		validate: validate,
		logger:   logger,
		tag2error: map[string]errors.UserProtoError{
			"http_url": createUserProtoError("http_url", false),
			"required": createUserProtoError("required", false),
			"gte":      createUserProtoError("gte", true),
			"lte":      createUserProtoError("lte", true),
			"max":      createUserProtoError("max", true),
			"min":      createUserProtoError("min", true),
		},
	}
}

// Register - регистрирует новые именованные функции валидации полей.
func (v *ValidatorAdapter) Register(tagName string, fn func(value string) bool) error {
	v.tag2error[tagName] = createUserProtoError(tagName, false)

	return v.validate.RegisterValidation(
		tagName,
		func(fl validator.FieldLevel) bool {
			return fn(fl.Field().String())
		},
	)
}

// Validate - возвращает результат валидации указанной структуру
// или ошибку с полями, в которых обнаружены проблемы.
func (v *ValidatorAdapter) Validate(ctx context.Context, structure any) error {
	err := v.validate.Struct(structure)

	// errors not found, OK
	if err == nil {
		return nil
	}

	errorList, ok := err.(validator.ValidationErrors) //nolint:errorlint
	if !ok {
		return errors.WrapInternalError(err, "err is not validator.ValidationErrors")
	}

	fields := make([]errors.CustomError, len(errorList))

	for i, errField := range errorList {
		fieldName := xstrings.TrimBeforeSep(errField.Namespace(), '.')
		fields[i] = errors.WithCustomCode(
			v.createUserError(ctx, fieldName, errField),
			fieldName,
		)

		v.logger.DebugFunc(
			ctx,
			func() string {
				return fmt.Sprintf(
					"{Field: %s, "+
						"StructNamespace: %s, "+
						"Tag: %s, "+
						"ActualTag: %s, "+
						"Kind: %v, "+
						"Type: %v, "+
						"Value: %v, "+
						"Param: %s}",
					fieldName,
					errField.StructNamespace(),
					errField.Tag(),
					errField.ActualTag(),
					errField.Kind(),
					errField.Type(),
					errField.Value(),
					errField.Param(),
				)
			},
			"validate.field", errField.Namespace(),
		)
	}

	return errors.CustomListError(fields)
}

func (v *ValidatorAdapter) createUserError(ctx context.Context, fieldName string, field validator.FieldError) error {
	tag := field.Tag()

	if tag == "" {
		return errInternalValidatorTagIsNotFound.New("field", fieldName) // TODO: нужна ли эта проверка?
	}

	args := make([]any, 0, 4)
	args = append(args, fieldName, field.Kind().String(), field.Value())

	if field.Param() != "" {
		args = append(args, field.Param())
	}

	if e, ok := v.tag2error[tag]; ok {
		return e.New(args...)
	}

	v.logger.Warn(ctx, "validator tag not registered", "tag", tag, "fieldName", fieldName)

	return createUserProtoError(tag, field.Param() != "").New(args...)
}

func createUserProtoError(tag string, withParam bool) errors.UserProtoError {
	if withParam {
		return errors.NewUserProto(
			validatorErrorID,
			validatorErrorPrefix+tag+validatorErrorPostfixWithParam,
		)
	}

	return errors.NewUserProto(
		validatorErrorID,
		validatorErrorPrefix+tag+validatorErrorPostfix,
	)
}
