package mrreq

import (
	"strings"
)

// ParseAccessToken - возвращает токен из заголовка Authorization.
// Если заголовка нет или он пустой, то вернётся пустое значение.
func ParseAccessToken(getter valueGetter) string {
	if token, ok := strings.CutPrefix(getter.Get("Authorization"), "Bearer "); ok {
		return strings.TrimSpace(token)
	}

	return ""
}
