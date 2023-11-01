package mrperms

import (
    "errors"
    "fmt"
    "strings"

    "github.com/mondegor/go-webcore/mrcore"
)

type (
    ClientSection struct {
        Caption string
        name string
        privilege string
        access *ModulesAccess
    }
)

func NewClientSection(name string, caption string, priv string, access *ModulesAccess) *ClientSection {
    if _, ok := access.privileges[priv]; !ok {
        mrcore.LogWarn("privilege '%s' is not registered in ModulesAccess.privileges, perhaps, it is not related to any role", priv)
    }

    return &ClientSection{
        Caption: caption,
        name: strings.Trim(name, "/"),
        privilege: priv,
        access: access,
    }
}

func (s *ClientSection) Path(actionPath string) string {
    return fmt.Sprintf("/%s/%s", s.name, strings.TrimLeft(actionPath, "/"))
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
