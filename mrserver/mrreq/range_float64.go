package mrreq

import (
	"net/http"

	"github.com/mondegor/go-webcore/mrtype"
)

// ParseRangeFloat64 - comment func.
func ParseRangeFloat64(r *http.Request, key string) (mrtype.RangeFloat64, error) {
	minValue, err := ParseFloat64(r, key+"-min", false)
	if err != nil {
		return mrtype.RangeFloat64{}, err
	}

	maxValue, err := ParseFloat64(r, key+"-max", false)
	if err != nil {
		return mrtype.RangeFloat64{}, err
	}

	if maxValue > 0 && minValue > maxValue { // change
		return mrtype.RangeFloat64{
			Min: maxValue,
			Max: minValue,
		}, nil
	}

	return mrtype.RangeFloat64{
		Min: minValue,
		Max: maxValue,
	}, nil
}
