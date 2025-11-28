package config

type (
	// ActionGroup - группа (секция) объединяющая несколько обработчиков.
	// Config версия mraccess.ActionGroup.
	ActionGroup struct {
		Name      string `yaml:"name"`
		BasePath  string `yaml:"base_path"`
		Privilege string `yaml:"privilege"`
	}
)
