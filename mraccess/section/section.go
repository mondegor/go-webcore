package section

import (
	"strings"
)

// TODO: возможно переименовать пакет, т.к. локальная переменная затеняет пакет

type (
	// RoutingSection - маршрутная секция с базовым путём и собственными правами.
	RoutingSection struct {
		name      string
		basePath  string
		privilege string
	}
)

// NewRoutingSection - создаёт объект RoutingSection.
func NewRoutingSection(
	name string,
	basePath string,
	privilege string,
) *RoutingSection {
	basePath = "/" + strings.Trim(basePath, "/")

	if basePath != "/" {
		basePath += "/"
	}

	return &RoutingSection{
		name:      name,
		basePath:  basePath,
		privilege: privilege,
	}
}

// Name - название секции.
func (s *RoutingSection) Name() string {
	return s.name
}

// BuildPath - возвращает путь маршрута вместе с базовым путем секции.
func (s *RoutingSection) BuildPath(routePath string) string {
	if routePath == "" {
		return ""
	}

	return s.basePath + strings.TrimLeft(routePath, "/")
}

// Privilege - возвращает название привилегии секции.
func (s *RoutingSection) Privilege() string {
	return s.privilege
}
