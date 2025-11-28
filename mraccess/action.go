package mraccess

type (
	// Action - обработчик с наделёнными ему правами доступа.
	Action struct {
		Name       string
		Privilege  string
		Permission string
	}
)
