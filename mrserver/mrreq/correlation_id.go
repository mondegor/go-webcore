package mrreq

import (
	"net/http"
	"regexp"
	"strings"
)

const (
	HeaderKeyCorrelationID = "X-Correlation-Id"

	// f7479171-83d2-4f64-84ac-892f8c0aaf48
	minLenCorrelationID = 16
	maxLenCorrelationID = 64
)

var regexpCorrelationID = regexp.MustCompile(`^[0-9a-fA-F][0-9a-fA-F-]+[0-9a-fA-F]$`)

func ParseCorrelationID(r *http.Request) (string, error) {
	value := strings.TrimSpace(r.Header.Get(HeaderKeyCorrelationID))

	if value == "" {
		return "", nil
	}

	if len(value) < minLenCorrelationID ||
		len(value) > maxLenCorrelationID ||
		!regexpCorrelationID.MatchString(value) {
		return "", FactoryErrHTTPRequestCorrelationID.New(value)
	}

	return value, nil
}
