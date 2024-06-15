package mrpath

type (
	// PathBuilder - дополняет и формирует указанный путь.
	PathBuilder interface {
		BuildPath(path string) string
	}
)
