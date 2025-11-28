package request

import (
	"net/http"
	"strings"
)

// AccessToken - возвращает токен из заголовка Authorization.
// Если заголовка нет или он пустой, то вернётся пустое значение.
func AccessToken(r *http.Request) string {
	if token, ok := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer "); ok {
		return strings.TrimSpace(token)
	}

	return ""
}
