package mraccess

import (
	"strings"

	"github.com/mondegor/go-sysmess/mrpath"
)

type (
	// ActionGroup - группа (секция) объединяющая несколько обработчиков.
	// Включает базовый путь к этим обработчикам и обладает привилегией доступа к ним.
	ActionGroup struct {
		// Name - уникальное имя группы действий.
		Name string

		// Privilege - привилегия, необходимая для доступа ко всем действиям в группе.
		Privilege string

		// BasePath - базовый путь для всех обработчиков в группе.
		BasePath mrpath.Builder
	}
)

// NewActionGroup - создаёт новую группу действий с указанным именем, привилегией и базовым путём.
// Базовый путь автоматически нормализуется (добавляется ведущий слэш, убираются дубликаты).
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
