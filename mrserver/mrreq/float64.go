package mrreq

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/mondegor/go-webcore/mrcore"
)

const (
	maxLenFloat64 = 64
)

// ParseFloat64 - comment func.
func ParseFloat64(r *http.Request, key string, required bool) (float64, error) {
	value := strings.TrimSpace(r.URL.Query().Get(key))

	if value == "" {
		if required {
			return 0, mrcore.ErrHttpRequestParamEmpty.New(key)
		}

		return 0, nil
	}

	if len(value) > maxLenFloat64 {
		return 0, mrcore.ErrHttpRequestParamLenMax.New(key, maxLenFloat64)
	}

	item, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, mrcore.ErrHttpRequestParseParam.Wrap(err, "Float64", value)
	}

	return item, nil
}
