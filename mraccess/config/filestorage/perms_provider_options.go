package filestorage

// Option - comment func type.
type (
	Option func(p *permsProvider)
)

// WithPrivileges - comment func.
func WithPrivileges(values []string) Option {
	return func(p *permsProvider) {
		p.privileges = values
	}
}

// WithPermissions - comment func.
func WithPermissions(values []string) Option {
	return func(p *permsProvider) {
		p.permissions = values
	}
}
