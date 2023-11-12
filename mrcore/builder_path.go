package mrcore

type (
	BuilderPath interface {
		FullPath(path string) string
	}
)
