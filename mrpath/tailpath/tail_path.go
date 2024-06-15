package tailpath

import (
	"strings"

	"github.com/mondegor/go-webcore/mrpath"
)

type (
	// Builder - comment struct.
	Builder struct {
		basePath string
	}
)

// Make sure the Image conforms with the mrpath.PathBuilder interface.
var _ mrpath.PathBuilder = (*Builder)(nil)

// New - создаёт объект NewBuilderPath.
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
