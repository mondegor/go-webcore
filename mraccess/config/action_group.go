package config

type (
	// ActionGroup - группа (секция) объединяющая несколько обработчиков.
	// Config-версия mraccess.ActionGroup для парсинга из YAML-файлов.
	ActionGroup struct {
		// Name - уникальное имя группы действий.
		Name string `yaml:"name"`

		// BasePath - базовый путь для всех обработчиков в группе.
		BasePath string `yaml:"base_path"`

		// Privilege - привилегия, необходимая для доступа ко всем действиям в группе.
		Privilege string `yaml:"privilege"`
	}
)
