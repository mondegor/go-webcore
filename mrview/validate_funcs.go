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

// NewValidateAND - создаёт составной валидатор, требующий выполнения ВСЕХ переданных условий.
// Если список conditions пуст, всегда возвращает false.
func NewValidateAND(conditions ...func(value string) bool) func(value string) bool {
	if len(conditions) == 0 {
		return func(_ string) bool {
			return false
		}
	}

	return func(value string) bool {
		for _, fn := range conditions {
			if !fn(value) {
				return false
			}
		}

		return true
	}
}

// NewValidateOR - создаёт составной валидатор, требующий выполнения ХОТЯ БЫ ОДНОГО условия.
// Если список conditions пуст, всегда возвращает false.
func NewValidateOR(conditions ...func(value string) bool) func(value string) bool {
	if len(conditions) == 0 {
		return func(_ string) bool {
			return false
		}
	}

	return func(value string) bool {
		for _, fn := range conditions {
			if fn(value) {
				return true
			}
		}

		return false
	}
}

// NewValidateInArray - создаёт валидатор, проверяющий наличие значения в списке допустимых.
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

// ValidateAnyNotSpaceSymbol - проверяет, что значение не содержит пробельных символов.
// Пустая строка также считается валидной.
func ValidateAnyNotSpaceSymbol(value string) bool {
	return regexpAnyNotSpaceSymbol.MatchString(value)
}

// ValidateVariable - проверяет, что значение является допустимым идентификатором (переменной).
// Допустимые символы: буквы (a-z, A-Z) и цифры, но начинаться должна с буквы.
// Примеры валидных значений: "myVar", "UserName", "count123".
func ValidateVariable(value string) bool {
	return regexpVariable.MatchString(value)
}

// ValidateName - проверяет, что значение является допустимым именем ресурса.
// Допустимые символы: буквы, цифры, "/", ".", "+", "-", "_".
// Не может начинаться или заканчиваться на специальные символы.
// Минимальная длина: 2 символа (первый и последний должны быть alphanumeric).
// Примеры валидных значений: "user_name", "my-app", "file.txt", "path/to/resource".
func ValidateName(value string) bool {
	return regexpName.MatchString(value)
}

// ValidateRewriteName - проверяет, что значение является допустимым человеко-читаемым именем (slug).
// Допустимые символы: строчные буквы (a-z), цифры и дефис.
// Не может начинаться или заканчиваться на дефис.
// Минимальная длина: 2 символа.
// Примеры валидных значений: "my-page", "user-profile", "item-123".
func ValidateRewriteName(value string) bool {
	return regexpRewriteName.MatchString(value)
}

// ValidatePassword - проверяет, что значение содержит допустимые символы для пароля.
// Допустимые символы: буквы, цифры и специальные символы (без пробелов).
func ValidatePassword(value string) bool {
	return regexpPassword.MatchString(value)
}

// ValidateDoubleSize - проверяет, что значение соответствует формату двумерного размера (ШxВ).
// Ожидаемый формат: "{ширина}x{высота}", где ширина и высота - целые числа.
// Примеры валидных значений: "100x200", "1920x1080".
func ValidateDoubleSize(value string) bool {
	return regexpDoubleSize.MatchString(value)
}

// ValidateTripleSize - проверяет, что значение соответствует формату трёхмерного размера (ШxВxГ).
// Ожидаемый формат: "{ширина}x{высота}x{глубина}", где все значения - целые числа.
// Примеры валидных значений: "100x200x300", "10x20x5".
func ValidateTripleSize(value string) bool {
	return regexpTripleSize.MatchString(value)
}
