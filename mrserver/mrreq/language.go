package mrreq

import (
	"net/http"

	"github.com/mondegor/go-sysmess/mrlang"
)

const (
	headerKeyAcceptLanguage = "Accept-Language"
)

// ParseLanguage - comment func.
func ParseLanguage(r *http.Request) []string {
	acceptLanguage := r.Header.Get(headerKeyAcceptLanguage)

	return mrlang.ParseAcceptLanguage(acceptLanguage)
}
