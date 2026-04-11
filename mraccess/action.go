package mraccess

type (
	// Action - обработчик с наделёнными ему правами доступа.
	Action struct {
		// Name - уникальное имя обработчика.
		Name string

		// Privilege - привилегия, необходимая для вызова обработчика.
		Privilege string

		// Permission - дополнительное разрешение для тонкой настройки доступа.
		Permission string
	}
)
