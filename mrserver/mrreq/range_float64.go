package mrreq

import (
	"github.com/mondegor/go-webcore/mrtype"
)

// ParseRangeFloat64 - возвращает RangeFloat64 из строковых параметров по указанному префиксу ключа.
func ParseRangeFloat64(getter valueGetter, prefixKey string) (mrtype.RangeFloat64, error) {
	minValue, err := ParseFloat64(getter, prefixKey+"-min", false)
	if err != nil {
		return mrtype.RangeFloat64{}, err
	}

	maxValue, err := ParseFloat64(getter, prefixKey+"-max", false)
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
