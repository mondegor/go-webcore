package mrperms

import (
	"strings"
)

type (
	// AppSection - comment struct.
	AppSection struct {
		caption   string
		rootPath  string
		privilege string
		secret    string
		audience  string
	}

	// AppSectionOptions - опции для создания AppSection.
	AppSectionOptions struct {
		Caption      string
		BasePath     string
		Privilege    string
		AuthSecret   string
		AuthAudience string
	}
)

// NewAppSection - создаёт объект AppSection.
func NewAppSection(opts AppSectionOptions) *AppSection {
	basePath := "/" + strings.Trim(opts.BasePath, "/")

	if basePath != "/" {
		basePath += "/"
	}

	return &AppSection{
		caption:   opts.Caption,
		rootPath:  basePath,
		privilege: opts.Privilege,
		secret:    opts.AuthSecret,
		audience:  opts.AuthAudience,
	}
}

// Caption - comment method.
func (s *AppSection) Caption() string {
	return s.caption
}

// BuildPath - comment method.
func (s *AppSection) BuildPath(methodPath string) string {
	if methodPath == "" {
		return ""
	}

	return s.rootPath + strings.TrimLeft(methodPath, "/")
}

// Privilege - comment method.
func (s *AppSection) Privilege() string {
	return s.privilege
}

// Secret - comment method.
func (s *AppSection) Secret() string {
	return s.secret
}

// Audience - comment method.
func (s *AppSection) Audience() string {
	return s.audience
}
