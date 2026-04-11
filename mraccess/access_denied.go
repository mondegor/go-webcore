package mraccess

type (
	// accessDenied - заглушка, используемая когда доступ явно запрещён.
	// Применяется как объект по умолчанию для неизвестных групп ролей.
	accessDenied struct{}
)

// HasPrivilege - всегда сообщает что привилегия отсутствует.
func (accessDenied) HasPrivilege(_ string) bool {
	return false
}

// HasPermission - всегда сообщает что разрешение отсутствует.
func (accessDenied) HasPermission(_ string) bool {
	return false
}
