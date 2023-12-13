package mrperms

import (
	"strings"

	"github.com/mondegor/go-webcore/mrcore"
)

type (
	ClientSection struct {
		caption   string
		rootPath  string
		privilege string
		secret    string
		audience  string
		access    *ModulesAccess
	}

	ClientSectionOptions struct {
		Caption      string
		RootPath     string
		Privilege    string
		AuthSecret   string
		AuthAudience string
		Access       *ModulesAccess
	}
)

func NewClientSection(opt ClientSectionOptions) *ClientSection {
	if _, ok := opt.Access.privileges[opt.Privilege]; !ok {
		mrcore.LogWarning(
			"privilege '%s' is not registered in ModulesAccess.privileges, perhaps, it is not related to any role",
			opt.Privilege,
		)
	}

	rootPath := "/" + strings.Trim(opt.RootPath, "/")

	if rootPath != "/" {
		rootPath += "/"
	}

	return &ClientSection{
		caption:   opt.Caption,
		rootPath:  rootPath,
		privilege: opt.Privilege,
		secret:    opt.AuthSecret,
		audience:  opt.AuthAudience,
		access:    opt.Access,
	}
}

func (s *ClientSection) Caption() string {
	return s.caption
}

func (s *ClientSection) Path(actionPath string) string {
	return s.rootPath + strings.TrimLeft(actionPath, "/")
}

func (s *ClientSection) MiddlewareWithPermission(name string, next mrcore.HttpHandlerFunc) mrcore.HttpHandlerFunc {
	return func(c mrcore.ClientContext) error {
		group := s.access.NewRoleGroup([]string{"administrators", "guests"}) // :TODO: брать у пользователя

		if !group.CheckPrivilege(s.privilege) && !group.CheckPermission(name) {
			if group.IsAuthorized() {
				return mrcore.FactoryErrHttpAccessForbidden.New()
			} else {
				return mrcore.FactoryErrHttpClientUnauthorized.New()
			}
		}

		return next(c)
	}
}
