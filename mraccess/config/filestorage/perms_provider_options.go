package filestorage

// Option - настройка объекта permsProvider.
type (
	Option func(o *permsProvider)
)

// WithPrivileges - comment func.
func WithPrivileges(values []string) Option {
	return func(o *permsProvider) {
		o.privileges = values
	}
}

// WithPermissions - comment func.
func WithPermissions(values []string) Option {
	return func(o *permsProvider) {
		o.permissions = values
	}
}
