package mrlib

import (
    "strings"
)

type (
    BuilderPath struct {
        basePath string
        postfix string
    }
)

func NewBuilderPath(basePath string) *BuilderPath {
    basePath = strings.TrimRight(basePath, "/")

    i := strings.Index(basePath, ":path:")

    if i > 0 {
        return &BuilderPath{
            basePath: basePath[0:i],
            postfix: basePath[i + 6:],
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
