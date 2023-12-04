package mrview

import (
	"context"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrctx"
)

// go get -u github.com/go-playground/validator/v10

const (
	validatorErrorPrefix = "validator_err"
)

type (
	validatorAdapter struct {
		validate *validator.Validate
	}
)

// Make sure the validatorAdapter conforms with the mrcore.Validator interface
var _ mrcore.Validator = (*validatorAdapter)(nil)

func NewValidator() *validatorAdapter {
	validate := validator.New()

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}

		return name
	})

	return &validatorAdapter{
		validate: validate,
	}
}

func (v *validatorAdapter) Register(tagName string, fn mrcore.ValidatorTagNameFunc) error {
	return v.validate.RegisterValidation(
		tagName,
		func(fl validator.FieldLevel) bool {
			return fn(fl.Field().String())
		},
	)
}

func (v *validatorAdapter) Validate(ctx context.Context, structure any) error {
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

func (v *validatorAdapter) createAppError(field validator.FieldError) *mrerr.AppError {
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
