package mrreq

import (
	"regexp"
	"strings"
)

const (
	minLenCorrelationID = 4
	maxLenCorrelationID = 64
)

var regexpCorrelationID = regexp.MustCompile(`^[0-9a-zA-Z][0-9a-zA-Z-]+[0-9a-zA-Z]$`)

// ParseCorrelationID - возвращает значение заголовка CorrelationID.
// Если заголовка нет или он пустой, то вернётся пустое значение.
func ParseCorrelationID(getter valueGetter) (string, error) {
	value := strings.TrimSpace(getter.Get(HeaderKeyCorrelationID))

	if value == "" {
		return "", nil
	}

	if len(value) < minLenCorrelationID ||
		len(value) > maxLenCorrelationID ||
		!regexpCorrelationID.MatchString(value) {
		return "", ErrHttpRequestCorrelationID.New(value)
	}

	return value, nil
}
