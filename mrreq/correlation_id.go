package mrreq

import (
	"net/http"
	"regexp"
	"strings"
)

const (
	// f7479171-83d2-4f64-84ac-892f8c0aaf48
	lenCorrelationID       = 36
	headerKeyCorrelationID = "CorrelationID"
)

var (
	regexpCorrelationID = regexp.MustCompile(`^[0-9a-fA-F][0-9a-fA-F-]{34}[0-9a-fA-F]$`)
)

func ParseCorrelationID(r *http.Request) (string, error) {
	value := strings.TrimSpace(r.Header.Get(headerKeyCorrelationID))

	if value == "" {
		return "", nil
	}

	if len(value) != lenCorrelationID ||
		!regexpCorrelationID.MatchString(value) {
		return "", FactoryErrHttpRequestCorrelationID.New(value)
	}

	return value, nil
}
