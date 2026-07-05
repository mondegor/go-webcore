package request

import (
	"net/http"
	"strings"
)

// AccessToken - извлекает токен доступа из заголовка Authorization.
// Ожидает формат: "Bearer <token>".
// Возвращает токен без ведущих/ведомых пробелов или пустую строку,
// если заголовок отсутствует или не имеет префикса "Bearer ".
func AccessToken(r *http.Request) string {
	if token, ok := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer "); ok {
		return strings.TrimSpace(token)
	}

	return ""
}
