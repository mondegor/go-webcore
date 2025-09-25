package groups

type (
	accessDenied struct{}
)

// CheckPrivilege - comment method.
func (accessDenied) CheckPrivilege(_ string) bool {
	return false
}

// CheckPermission - comment method.
func (accessDenied) CheckPermission(_ string) bool {
	return false
}
