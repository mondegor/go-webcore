package mrlib

import (
	"strings"
)

const (
	placeholderPath = "{{path}}"
)

type (
	BuilderPath struct {
		basePath string
		postfix  string
	}
)

// NewBuilderPath - sample /dir/{{path}}/postfix -> /dir/real-value/postfix
func NewBuilderPath(basePath string) *BuilderPath {
	return NewBuilderPathWithPlaceholder(basePath, placeholderPath)
}

func NewBuilderPathWithPlaceholder(basePath, placeholder string) *BuilderPath {
	basePath = strings.TrimRight(basePath, "/")

	if i := strings.Index(basePath, placeholder); i > 0 {
		return &BuilderPath{
			basePath: strings.TrimRight(basePath[0:i], "/"),
			postfix:  strings.TrimLeft(basePath[i+len(placeholder):], "/"),
		}
	}

	return &BuilderPath{
		basePath: basePath,
	}
}

func (p *BuilderPath) FullPath(path string) string {
	if path == "" {
		return ""
	}

	if p.postfix == "" {
		return strings.Join(
			[]string{
				p.basePath,
				strings.TrimLeft(path, "/"),
			},
			"/",
		)
	}

	return strings.Join(
		[]string{
			p.basePath,
			strings.Trim(path, "/"),
			p.postfix,
		},
		"/",
	)
}
