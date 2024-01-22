package mrplayvalidator

import (
	"context"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrctx"
	"github.com/mondegor/go-webcore/mrview"
)

// go get -u github.com/go-playground/validator/v10

const (
	validatorErrorPrefix = "validator_err"
)

type (
	ValidatorAdapter struct {
		validate *validator.Validate
	}
)

// Make sure the ValidatorAdapter conforms with the mrview.Validator interface
var _ mrview.Validator = (*ValidatorAdapter)(nil)

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

func (v *ValidatorAdapter) Register(tagName string, fn mrview.ValidatorTagNameFunc) error {
	return v.validate.RegisterValidation(
		tagName,
		func(fl validator.FieldLevel) bool {
			return fn(fl.Field().String())
		},
	)
}

func (v *ValidatorAdapter) Validate(ctx context.Context, structure any) error {
	err := v.validate.Struct(structure)

	// error not found, OK
	if err == nil {
		return nil
	}

	errors, ok := err.(validator.ValidationErrors)

	if !ok {
		return mrcore.FactoryErrInternal.Wrap(err)
	}

	fields := make([]*mrerr.FieldError, len(errors))
	logger := mrctx.Logger(ctx)

	for i, errField := range errors {
		fields[i] = mrerr.NewFieldErrorAppError(
			errField.Field(),
			v.createAppError(errField),
		)

		logger.Debug(
			"{Namespace: %s, "+
				"Field: %s, "+
				"StructNamespace: %s, "+
				"StructField: %s, "+
				"Tag: %s, "+
				"ActualTag: %s, "+
				"Kind: %v, "+
				"Type: %v, "+
				"Value: %v, "+
				"Param: %s}",
			errField.Namespace(),
			errField.Field(),
			errField.StructNamespace(),
			errField.StructField(),
			errField.Tag(),
			errField.ActualTag(),
			errField.Kind(),
			errField.Type(),
			errField.Value(),
			errField.Param(),
		)
	}

	return mrerr.FieldErrorList(fields)
}

func (v *ValidatorAdapter) createAppError(field validator.FieldError) *mrerr.AppError {
	tag := field.Tag()

	if len(tag) == 0 {
		return mrerr.New(validatorErrorPrefix, validatorErrorPrefix)
	}

	id := validatorErrorPrefix + "_" + tag
	message := id + ": name={{ .name }}, type={{ .type }}, value={{ .value }}"
	param := field.Param()

	if param != "" {
		return mrerr.New(
			id,
			message+", param={{ .param }}",
			field.Field(),
			field.Kind().String(),
			field.Value(),
			param,
		)
	}

	return mrerr.New(
		id,
		message,
		field.Field(),
		field.Kind().String(),
		field.Value(),
	)
}
