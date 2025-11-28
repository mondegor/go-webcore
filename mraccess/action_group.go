package mraccess

import (
	"strings"

	"github.com/mondegor/go-sysmess/mrpath"
)

type (
	// ActionGroup - группа (секция) объединяющая несколько обработчиков.
	// Включает базовый путь к этим обработчикам и обладает привилегией доступа к ним.
	ActionGroup struct {
		Name      string
		Privilege string
		BasePath  mrpath.Builder
	}
)

// NewActionGroup - создаёт объект ActionGroup.
func NewActionGroup(
	name string,
	privilege string,
	basePath string,
) *ActionGroup {
	return &ActionGroup{
		Name:      name,
		Privilege: privilege,
		BasePath:  mrpath.NewTail("/" + strings.Trim(basePath, "/")),
	}
}
