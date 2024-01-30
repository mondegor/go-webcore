package mrperms

import (
	"strings"
)

type (
	appSection struct {
		caption   string
		rootPath  string
		privilege string
		secret    string
		audience  string
	}

	AppSectionOptions struct {
		Caption      string
		RootPath     string
		Privilege    string
		AuthSecret   string
		AuthAudience string
	}
)

func NewAppSection(opts AppSectionOptions) AppSection {
	rootPath := "/" + strings.Trim(opts.RootPath, "/")

	if rootPath != "/" {
		rootPath += "/"
	}

	return &appSection{
		caption:   opts.Caption,
		rootPath:  rootPath,
		privilege: opts.Privilege,
		secret:    opts.AuthSecret,
		audience:  opts.AuthAudience,
	}
}

func (s *appSection) Caption() string {
	return s.caption
}

func (s *appSection) Path(actionPath string) string {
	return s.rootPath + strings.TrimLeft(actionPath, "/")
}

func (s *appSection) Privilege() string {
	return s.privilege
}

func (s *appSection) Secret() string {
	return s.secret
}

func (s *appSection) Audience() string {
	return s.audience
}
