package mrview

type (
	// Tag - представляет тег валидации с именем и ассоциированной функцией проверки.
	// Используется для регистрации кастомных правил валидации в Validator.
	Tag struct {
		// Name - уникальное имя тега для использования в struct-аннотациях.
		Name string

		// ValidateFunc - функция, принимающая строковое значение и возвращающая результат валидации.
		ValidateFunc func(value string) bool
	}
)

// TagArticle - создаёт тег для валидации слова-идентификатора.
// Примеры валидных значений: "article123", "my-article", "123".
func TagArticle() Tag {
	return Tag{
		Name:         "tag_article",
		ValidateFunc: ValidateAnyNotSpaceSymbol,
	}
}

// TagVariable - создаёт тег для валидации имени переменной.
// Примеры валидных значений: "myVar", "UserName", "count123".
func TagVariable() Tag {
	return Tag{
		Name:         "tag_variable",
		ValidateFunc: ValidateVariable,
	}
}

// TagName - создаёт тег для валидации имени.
// Примеры валидных значений: "user_name", "my-app", "file.txt", "path/to/resource".
func TagName() Tag {
	return Tag{
		Name:         "tag_name",
		ValidateFunc: ValidateName,
	}
}

// TagRewriteName - создаёт тег для валидации человеко-читаемого имени.
// Примеры валидных значений: "my-page", "user-profile", "item-123".
func TagRewriteName() Tag {
	return Tag{
		Name:         "tag_rewrite_name",
		ValidateFunc: ValidateRewriteName,
	}
}

// TagPassword - создаёт тег для валидации пароля.
func TagPassword() Tag {
	return Tag{
		Name:         "tag_password",
		ValidateFunc: ValidatePassword,
	}
}

// TagDoubleSize - создаёт тег для валидации двумерного размера (ШxВ).
// Примеры валидных значений: "100x200", "1920x1080".
func TagDoubleSize() Tag {
	return Tag{
		Name:         "tag_double_size",
		ValidateFunc: ValidateDoubleSize,
	}
}

// TagTripleSize - создаёт тег для валидации трёхмерного размера (ШxВxГ).
// Примеры валидных значений: "100x200x300", "10x20x5".
func TagTripleSize() Tag {
	return Tag{
		Name:         "tag_triple_size",
		ValidateFunc: ValidateTripleSize,
	}
}
