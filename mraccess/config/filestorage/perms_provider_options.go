package filestorage

// Option - настройка объекта permsProvider.
type (
	Option func(o *permsProvider)
)

// WithPrivileges - устанавливает список привилегий для поставщика прав доступа.
func WithPrivileges(values []string) Option {
	return func(o *permsProvider) {
		o.privileges = values
	}
}

// WithPermissions - устанавливает список разрешений для поставщика прав доступа.
func WithPermissions(values []string) Option {
	return func(o *permsProvider) {
		o.permissions = values
	}
}
