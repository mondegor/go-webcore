package mrreq

import (
	"net/http"
	"strings"
	"time"

	"github.com/mondegor/go-webcore/mrcore"
)

const (
	maxLenDateTime = 64
)

// ParseDateTime - comment func.
func ParseDateTime(r *http.Request, key string, required bool) (time.Time, error) {
	value := strings.TrimSpace(r.URL.Query().Get(key))

	if value == "" {
		if required {
			return time.Time{}, mrcore.ErrHttpRequestParamEmpty.New(key)
		}

		return time.Time{}, nil
	}

	if len(value) > maxLenDateTime {
		return time.Time{}, mrcore.ErrHttpRequestParamLenMax.New(key, maxLenDateTime)
	}

	item, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return time.Time{}, mrcore.ErrHttpRequestParseParam.Wrap(err, key, "DateTime", value)
	}

	return item, nil
}
