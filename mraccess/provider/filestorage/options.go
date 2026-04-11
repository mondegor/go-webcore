package filestorage

type (
	// Option - функциональная настройка объекта PermsProvider.
	Option func(o *options)

	options struct {
		allowedPrivileges  []string
		allowedPermissions []string
	}
)

// WithAllowedPrivileges - устанавливает список допустимых привилегий для поставщика прав доступа.
func WithAllowedPrivileges(values []string) Option {
	return func(o *options) {
		o.allowedPrivileges = values
	}
}

// WithAllowedPermissions - устанавливает список допустимых разрешений для поставщика прав доступа.
func WithAllowedPermissions(values []string) Option {
	return func(o *options) {
		o.allowedPermissions = values
	}
}
