package mrreq

import (
	"net/http"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrtype"
)

const (
	maxValuePageSize = 1000000
)

func ParsePageParams(r *http.Request, keyIndex, keySize string) (mrtype.PageParams, error) {
	params := mrtype.PageParams{}
	index, err := ParseInt64(r, keyIndex, false)

	if index < 0 {
		index = 0
	}

	if err != nil {
		return params, err
	}

	size, err := ParseInt64(r, keySize, false)

	if err != nil {
		return params, err
	}

	if size < 0 {
		size = 0
	}

	if size > maxValuePageSize {
		return params, mrcore.FactoryErrHttpRequestParamMax.New(keySize, maxValuePageSize)
	}

	params.Index = uint64(index)
	params.Size = uint64(size)

	return params, nil
}
