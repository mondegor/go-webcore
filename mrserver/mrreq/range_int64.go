package mrreq

import (
	"github.com/mondegor/go-webcore/mrtype"
)

// ParseRangeInt64 - возвращает RangeInt64 из строковых параметров по указанному префиксу ключа.
func ParseRangeInt64(getter valueGetter, prefixKey string) (mrtype.RangeInt64, error) {
	minValue, err := ParseInt64(getter, prefixKey+"-min", false)
	if err != nil {
		return mrtype.RangeInt64{}, err
	}

	maxValue, err := ParseInt64(getter, prefixKey+"-max", false)
	if err != nil {
		return mrtype.RangeInt64{}, err
	}

	if maxValue > 0 && minValue > maxValue { // change
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
