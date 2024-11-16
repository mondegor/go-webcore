package mrplayvalidator

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/mondegor/go-sysmess/mrerr"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrlib"
	"github.com/mondegor/go-webcore/mrlog"
)

// go get -u github.com/go-playground/validator/v10

const (
	validatorErrorPrefix = "validator.err"
)

type (
	// ValidatorAdapter - адаптер валидатора структур и их полей на базе тегов.
	ValidatorAdapter struct {
		validate *validator.Validate
	}
)

var errValidatorTagIsNotFound = mrerr.NewProto(
	validatorErrorPrefix, mrerr.ErrorKindUser, "validator error: tag is empty").New()

// New - создаёт объект ValidatorAdapter.
func New() *ValidatorAdapter {
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
	}
}

// Register - регистрирует новые именованные функции валидации полей.
func (v *ValidatorAdapter) Register(tagName string, fn func(value string) bool) error {
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
		return mrcore.ErrInternal.Wrap(err)
	}

	fields := make(mrerr.CustomErrors, len(errorList))
	logger := mrlog.Ctx(ctx)

	for i, errField := range errorList {
		fieldName := mrlib.CutBefore(errField.Namespace(), '.')
		fields[i] = mrerr.NewCustomError(
			fieldName,
			v.createAppError(fieldName, errField),
		)

		logger.Debug().Str("validate.field", errField.Namespace()).MsgFunc(
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
		)
	}

	return fields
}

func (v *ValidatorAdapter) createAppError(fieldName string, field validator.FieldError) *mrerr.AppError {
	if len(field.Tag()) == 0 {
		return errValidatorTagIsNotFound
	}

	// TODO: шаблон динамической ошибки можно кэшировать

	// здесь передаются все атрибуты поля, которые можно вывести пользователю,
	id := validatorErrorPrefix + "_" + field.Tag()
	message := id + ": name={{ .name }}, type={{ .type }}, value={{ .value }}"
	args := [4]any{fieldName, field.Kind().String(), field.Value()}

	if field.Param() != "" {
		args[3] = field.Param()

		return mrerr.NewProto(
			id,
			mrerr.ErrorKindUser,
			message+", param={{ .param }}",
		).New(args[:]...)
	}

	return mrerr.NewProto(
		id,
		mrerr.ErrorKindUser,
		message,
	).New(args[:3]...)
}
