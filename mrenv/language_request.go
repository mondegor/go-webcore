package mrenv

import (
    "net/http"

    "github.com/mondegor/go-sysmess/mrlang"
)

func LanguageFromRequest(r *http.Request) []string {
    acceptLanguage := r.Header.Get("Accept-Language")

    return mrlang.ParseAcceptLanguage(acceptLanguage)
}
