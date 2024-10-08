package mrreq

import (
	"github.com/mondegor/go-sysmess/mrlang"
)

// ParseLanguage - возвращает список языков из заголовка Accept-Language.
// Если заголовка нет или он пустой, то возвращается пустой массив.
func ParseLanguage(getter valueGetter) []string {
	acceptLanguage := getter.Get(HeaderKeyAcceptLanguage)

	return mrlang.ParseAcceptLanguage(acceptLanguage)
}
