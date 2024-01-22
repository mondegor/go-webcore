package mrreq

import (
	"net/http"

	"github.com/mondegor/go-webcore/mrtype"
)

func ParseRangeInt64(r *http.Request, key string) (mrtype.RangeInt64, error) {
	minValue, err := ParseInt64(r, key+"-min", false)

	if err != nil {
		return mrtype.RangeInt64{}, err
	}

	maxValue, err := ParseInt64(r, key+"-max", false)

	if err != nil {
		return mrtype.RangeInt64{}, err
	}

	if minValue > maxValue { // change
		return mrtype.RangeInt64{
			Min: maxValue,
			Max: minValue,
		}, nil
	}

	return mrtype.RangeInt64{
		Min: minValue,
		Max: maxValue,
	}, nil
}
