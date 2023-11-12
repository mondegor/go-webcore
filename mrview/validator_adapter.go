package mrview

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrctx"
)

// go get -u github.com/go-playground/validator/v10

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

	return &validatorAdapter {
		validate: validate,
	}
}

func (v *validatorAdapter) Register(tagName string, fn mrcore.ValidatorTagNameFunc) error {
	return v.validate.RegisterValidation(
		tagName,
		func (fl validator.FieldLevel) bool {
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

	errorList := mrerr.NewList()
	logger := mrctx.Logger(ctx)

	for _, errField := range errors {
		errorList.Add(errField.Field(), v.createAppError(errField))

		logger.Debug(
			"Namespace: %s\n"+
				"Field: %s\n"+
				"StructNamespace: %s\n"+
				"StructField: %s\n"+
				"Tag: %s\n"+
				"ActualTag: %s\n"+
				"Kind: %v\n"+
				"Type: %v\n"+
				"Value: %v\n"+
				"Param: %s",
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

	return errorList
}

func (v *validatorAdapter) createAppError(field validator.FieldError) *mrerr.AppError {
	id := []byte("errValidation")
	tag := []byte(field.Tag())

	if len(tag) == 0 {
		return mrerr.New(string(id), string(id))
	}

	tag[0] -= 32 // to uppercase first char
	id = append(id, tag...)

	return mrerr.New(
		string(id),
		fmt.Sprintf("%s: value='{{ .value }}'", id),
		field.Value(),
	)
}
