package mrenv

import "net/http"

func PlatformFromRequest(r *http.Request) (string, error) {
    value := r.Header.Get("Platform")

    if value == "" {
        return "", nil
    }

    if value == PlatformWeb {
        return PlatformWeb, nil
    }

    if value == PlatformMobile {
        return PlatformMobile, nil
    }

    return "", factoryErrHttpRequestPlatformValue.New(value)
}
