package mrperms

import (
	"errors"
	"fmt"
	"strings"

	"github.com/mondegor/go-webcore/mrcore"
)

type (
	ClientSection struct {
		caption   string
		rootPath  string
		privilege string
		access    *ModulesAccess
	}
)

func NewClientSection(caption, rootPath, priv string, access *ModulesAccess) *ClientSection {
	if _, ok := access.privileges[priv]; !ok {
		mrcore.LogWarning(
			"privilege '%s' is not registered in ModulesAccess.privileges, perhaps, it is not related to any role",
			priv,
		)
	}

	return &ClientSection{
		caption:   caption,
		rootPath:  strings.Trim(rootPath, "/"),
		privilege: priv,
		access:    access,
	}
}

func (s *ClientSection) Caption() string {
	return s.caption
}

func (s *ClientSection) Path(actionPath string) string {
	return fmt.Sprintf("/%s/%s", s.rootPath, strings.TrimLeft(actionPath, "/"))
}

func (s *ClientSection) MiddlewareWithPermission(name string, next mrcore.HttpHandlerFunc) mrcore.HttpHandlerFunc {
	return func(c mrcore.ClientData) error {
		group := s.access.NewRoleGroup([]string{"administrators", "guests"}) // :TODO: брать у пользователя

		if !group.CheckPrivilege(s.privilege) {
			return errors.New("403") // :TODO: превратить в системную ошибку
		}

		if !group.CheckPermission(name) {
			return errors.New("403") // :TODO: превратить в системную ошибку
		}

		return next(c)
	}
}
