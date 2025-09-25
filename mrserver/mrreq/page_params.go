package mrreq

import (
	"math"

	"github.com/mondegor/go-sysmess/mrerr/mr"

	"github.com/mondegor/go-webcore/mrtype"
)

// ParsePageParams - возвращает PageParams из строковых параметров по указанным ключам.
func ParsePageParams(getter valueGetter, keyIndex, keySize string) (mrtype.PageParams, error) {
	index, err := ParseUint64(getter, keyIndex, false)
	if err != nil {
		return mrtype.PageParams{}, err
	}

	size, err := ParseUint64(getter, keySize, false)
	if err != nil {
		return mrtype.PageParams{}, err
	}

	if size > math.MaxUint32 {
		return mrtype.PageParams{}, mr.ErrHttpRequestParamMax.New(keySize, math.MaxUint32)
	}

	if index > size {
		return mrtype.PageParams{}, mr.ErrHttpRequestParamMax.New(keyIndex, size)
	}

	return mrtype.PageParams{
		Index: index,
		Size:  size,
	}, nil
}
