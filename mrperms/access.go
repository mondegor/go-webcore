package mrperms

type (
	AppSection interface {
		Caption() string
		Path(actionPath string) string
		Privilege() string
		Secret() string
		Audience() string
	}

	AccessControl interface {
		NewAccessRights(members ...string) AccessRights
		HasPrivilege(name string) bool
		HasPermission(name string) bool
	}

	AccessRights interface {
		IsGuestAccess() bool
		CheckPrivilege(name string) bool
		CheckPermission(name string) bool
	}
)
