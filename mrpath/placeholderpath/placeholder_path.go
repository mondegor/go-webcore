package placeholderpath

import (
	"fmt"
	"strings"

	"github.com/mondegor/go-webcore/mrcore"
)

const (
	Placeholder = "{{path}}" // Placeholder - плейсхолдер в пути по умолчанию
)

type (
	// Builder - comment struct.
	Builder struct {
		basePath string
		postfix  string
	}
)

// New - создаёт объект Builder.
// sample: /dir/{{path}}/postfix -> /dir/real-value/postfix
func New(basePath, placeholder string) (*Builder, error) {
	if i := strings.Index(basePath, placeholder); i > 0 {
		return &Builder{
			basePath: basePath[0:i],
			postfix:  basePath[i+len(placeholder):],
		}, nil
	}

	return nil, mrcore.ErrInternal.Wrap(fmt.Errorf("placeholder '%s' is not found in path '%s'", placeholder, basePath))
}

// BuildPath - comment method.
func (p *Builder) BuildPath(path string) string {
	if path == "" {
		return ""
	}

	return p.basePath + strings.Trim(path, "/") + p.postfix
}
