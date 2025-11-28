package mraccess

type (
	accessDenied struct{}
)

// HasPrivilege - всегда возвращает false.
func (accessDenied) HasPrivilege(_ string) bool {
	return false
}

// HasPermission - всегда возвращает false.
func (accessDenied) HasPermission(_ string) bool {
	return false
}
