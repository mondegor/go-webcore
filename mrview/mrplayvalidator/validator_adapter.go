package mrplayvalidator

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/mondegor/go-core/errors"
	"github.com/mondegor/go-core/mrlog"
	"github.com/mondegor/go-core/util/xstrings"
)

// go get -u github.com/go-playground/validator/v10

const (
	// validatorErrorPrefix - префикс для идентификатора ошибки валидации.
	validatorErrorPrefix = "Validator_"

	// validatorErrorPostfix - шаблон сообщения об ошибке без параметра.
	// Плейсхолдеры: {Name} - имя поля, {Type} - тип, {Value} - значение.
	validatorErrorPostfix = ": {Caption}, {Type}, {Value}"

	// validatorErrorPostfixWithParam - шаблон сообщения об ошибке с параметром.
	// Дополнительный плейсхолдер: {Param} - значение параметра валидации.
	validatorErrorPostfixWithParam = validatorErrorPostfix + ", {Param}"

	// validatorErrorID - универсальный идентификатор ошибки валидации.
	validatorErrorID = "ValidateError"
)

type (
	// ValidatorAdapter - адаптер валидатора структур и их полей на базе тегов.
	// Использует библиотеку go-playground/validator для проверки структур.
	//
	// Особенности:
	//   - Имена полей в ошибках берутся из тега `json` структуры;
	//   - Поддержка регистрации пользовательских тегов валидации;
	//   - Преобразование ошибок в формат UserProtoError с идентификаторами;
	ValidatorAdapter struct {
		validate  *validator.Validate
		logger    mrlog.Logger
		tag2error map[string]errors.UserProtoError
	}
)

var errInternalValidatorTagIsNotFound = errors.NewInternalProto("validator error: tag is empty")

// New - создаёт и настраивает адаптер валидатора на базе go-playground/validator.
// Регистрируются шаблоны ошибок для популярных тегов:
//
//	http_url, required, gte, lte, max, min
//
// Возвращает настроенный экземпляр ValidatorAdapter, готовый к использованию.
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

// Register - регистрирует пользовательский тег валидации полей.
// Параметры:
//   - tagName - имя тега для использования в struct-аннотациях (например: `validate:"mytag"`);
//   - fn - функция валидации, принимающая строковое значение поля и возвращающая bool;
//
// При регистрации автоматически создаётся шаблон ошибки без параметра.
func (v *ValidatorAdapter) Register(tagName string, fn func(value string) bool) error {
	v.tag2error[tagName] = createUserProtoError(tagName, false)

	return v.validate.RegisterValidation(
		tagName,
		func(fl validator.FieldLevel) bool {
			return fn(fl.Field().String())
		},
	)
}

// Validate - выполняет валидацию указанной структуры по тегам.
// Возвращает CustomListError со списком ошибок для каждого проблемного поля.
// Каждая ошибка содержит:
//   - CustomCode - имя проблемного поля (из тега json);
//   - UserProtoError - код и описание ошибки в формате {Name}, {Type}, {Value}, {Param};
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
