package mrlib

import (
	"strings"
)

const (
	placeholderPath = "{{path}}"
)

type (
	BuilderPath interface {
		FullPath(path string) string
	}

	builderPath struct {
		basePath string
		postfix  string
	}
)

// NewBuilderPath - sample /dir/{{path}}/postfix -> /dir/real-value/postfix
func NewBuilderPath(basePath string) BuilderPath {
	return NewBuilderPathWithPlaceholder(basePath, placeholderPath)
}

func NewBuilderPathWithPlaceholder(basePath, placeholder string) BuilderPath {
	basePath = strings.TrimRight(basePath, "/")

	if i := strings.Index(basePath, placeholder); i > 0 {
		return &builderPath{
			basePath: strings.TrimRight(basePath[0:i], "/"),
			postfix:  strings.TrimLeft(basePath[i+len(placeholder):], "/"),
		}
	}

	return &builderPath{
		basePath: basePath,
	}
}

func (p *builderPath) FullPath(path string) string {
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
