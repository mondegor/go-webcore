package mrperms

import (
	"strings"
)

type (
	AppSection struct {
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

func NewAppSection(opt AppSectionOptions) *AppSection {
	rootPath := "/" + strings.Trim(opt.RootPath, "/")

	if rootPath != "/" {
		rootPath += "/"
	}

	return &AppSection{
		caption:   opt.Caption,
		rootPath:  rootPath,
		privilege: opt.Privilege,
		secret:    opt.AuthSecret,
		audience:  opt.AuthAudience,
	}
}

func (s *AppSection) Caption() string {
	return s.caption
}

func (s *AppSection) Path(actionPath string) string {
	return s.rootPath + strings.TrimLeft(actionPath, "/")
}

func (s *AppSection) Privilege() string {
	return s.privilege
}

func (s *AppSection) Secret() string {
	return s.secret
}

func (s *AppSection) Audience() string {
	return s.audience
}
