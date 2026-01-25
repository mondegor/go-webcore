package request

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/mondegor/go-webcore/mrserver"
)

const (
	minLenCorrelationID = 4
	maxLenCorrelationID = 64
)

var (
	regexpCorrelationID             = regexp.MustCompile(`^[0-9a-zA-Z][0-9a-zA-Z-]+[0-9a-zA-Z]$`)
	errHeaderContainsIncorrectValue = fmt.Errorf("header %s contains incorrect value", mrserver.HeaderKeyCorrelationID)
)

// CorrelationID - возвращает значение заголовка CorrelationID.
// Если заголовка нет или он пустой, то вернётся пустое значение.
func CorrelationID(r *http.Request) (string, error) {
	value := strings.TrimSpace(r.Header.Get(mrserver.HeaderKeyCorrelationID))

	if value == "" {
		return "", nil
	}

	if len(value) < minLenCorrelationID ||
		len(value) > maxLenCorrelationID ||
		!regexpCorrelationID.MatchString(value) {
		return "", errHeaderContainsIncorrectValue
	}

	return value, nil
}
