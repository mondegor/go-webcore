package mrreq

import (
    "net/http"

    "github.com/mondegor/go-webcore/mrtype"
)

func ParseRangeInt64(r *http.Request, key string) (mrtype.RangeInt64, error) {
    min, err := ParseInt64(r, key + "-min", false)

    if err != nil {
        return mrtype.RangeInt64{}, err
    }

    max, err := ParseInt64(r, key + "-max", false)

    if err != nil {
        return mrtype.RangeInt64{}, err
    }

    if min > max { // change
        return mrtype.RangeInt64{
            Min: max,
            Max: min,
        }, nil
    }

    return mrtype.RangeInt64{
        Min: min,
        Max: max,
    }, nil
}
