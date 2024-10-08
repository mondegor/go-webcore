package mrreq

type (
	valueGetter interface {
		Get(string) string
	}
)
