package mrview

import (
	"regexp"
)

var (
	regexpAnyNotSpaceSymbol = regexp.MustCompile(`^\S+$`)
	regexpVariable          = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9]*$`)
	regexpName              = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9/_.+-]*[a-zA-Z0-9]$`)
	regexpRewriteName       = regexp.MustCompile(`^[a-z][a-z0-9-]*[a-z0-9]$`)
	regexpPassword          = regexp.MustCompile(`^[a-zA-Z0-9!"#$%&'()*+,\-./:;<=>?@\[\\\]^_{|}~]+$`)
	regexpDoubleSize        = regexp.MustCompile(`^[0-9]+x[0-9]+$`)
	regexpTripleSize        = regexp.MustCompile(`^[0-9]+x[0-9]+x[0-9]+$`)
)

// NewValidateAND - создаёт составной валидатор, требующий выполнения всех условий.
func NewValidateAND(values ...func(value string) bool) func(value string) bool {
	if len(values) == 0 {
		return func(_ string) bool {
			return false
		}
	}

	return func(value string) bool {
		for _, fn := range values {
			if !fn(value) {
				return false
			}
		}

		return true
	}
}

// NewValidateOR - создаёт составной валидатор, требующий выполнения хотя бы одного условия.
func NewValidateOR(values ...func(value string) bool) func(value string) bool {
	return func(value string) bool {
		for _, fn := range values {
			if fn(value) {
				return true
			}
		}

		return false
	}
}

// NewValidateInArray - создаёт валидатор, проверяющий наличие значения в массиве.
func NewValidateInArray(items []string) func(value string) bool {
	return func(value string) bool {
		for _, item := range items {
			if item == value {
				return true
			}
		}

		return false
	}
}

// ValidateAnyNotSpaceSymbol - возвращает, не содержит ли значение пробельных символов.
func ValidateAnyNotSpaceSymbol(value string) bool {
	return regexpAnyNotSpaceSymbol.MatchString(value)
}

// ValidateVariable - проверяет, что значение является допустимым идентификатором.
func ValidateVariable(value string) bool {
	return regexpVariable.MatchString(value)
}

// ValidateName - проверяет, что значение является допустимым именем.
func ValidateName(value string) bool {
	return regexpName.MatchString(value)
}

// ValidateRewriteName - проверяет, что значение является допустимым человеко-понятным именем.
func ValidateRewriteName(value string) bool {
	return regexpRewriteName.MatchString(value)
}

// ValidatePassword - проверяет, что значение содержит допустимые символы для пароля.
func ValidatePassword(value string) bool {
	return regexpPassword.MatchString(value)
}

// ValidateDoubleSize - проверяет, что значение соответствует формату двойного размера (ШxВ).
func ValidateDoubleSize(value string) bool {
	return regexpDoubleSize.MatchString(value)
}

// ValidateTripleSize - проверяет, что значение соответствует формату тройного размера (ШxВxГ).
func ValidateTripleSize(value string) bool {
	return regexpTripleSize.MatchString(value)
}
