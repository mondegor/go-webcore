package mrreq

import (
    "net/http"
    "regexp"
)

const (
    // f7479171-83d2-4f64-84ac-892f8c0aaf48
	correlationIdLen= 36
)

var (
	regexpCorrelationId = regexp.MustCompile(`^[0-9a-fA-F][0-9a-fA-F-]{34}[0-9a-fA-F]$`)
)

func CorrelationId(r *http.Request) (string, error) {
    value := r.Header.Get("CorrelationID")

    if value == "" {
        return "", nil
    }

    if len(value) != correlationIdLen ||
       !regexpCorrelationId.MatchString(value) {
        return "", factoryErrHttpRequestCorrelationID.New(value)
    }

    return value, nil
}
