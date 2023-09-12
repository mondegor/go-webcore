package mrreq

import (
    "net/http"

    "github.com/mondegor/go-sysmess/mrlang"
)

func Language(r *http.Request) []string {
    acceptLanguage := r.Header.Get("Accept-Language")

    return mrlang.ParseAcceptLanguage(acceptLanguage)
}
