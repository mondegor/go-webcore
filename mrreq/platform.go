package mrreq

import (
    "net/http"

    "github.com/mondegor/go-webcore/mrcore"
)

func Platform(r *http.Request) (string, error) {
    value := r.Header.Get("Platform")

    if value == "" {
        return "", nil
    }

    if value == mrcore.PlatformWeb {
        return mrcore.PlatformWeb, nil
    }

    if value == mrcore.PlatformMobile {
        return mrcore.PlatformMobile, nil
    }

    return "", FactoryErrHttpRequestPlatformValue.New(value)
}
