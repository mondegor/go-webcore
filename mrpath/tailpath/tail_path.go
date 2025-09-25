package tailpath

import (
	"strings"
)

type (
	// Builder - comment struct.
	Builder struct {
		basePath string
	}
)

// New - создаёт объект Builder.
// sample: /base/dir/ -> /base/dir/path
func New(basePath string) *Builder {
	return &Builder{
		basePath: strings.TrimRight(basePath, "/") + "/",
	}
}

// BuildPath - comment method.
func (p *Builder) BuildPath(path string) string {
	if path == "" {
		return ""
	}

	return p.basePath + strings.TrimLeft(path, "/")
}
