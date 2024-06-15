package mrreq

import (
	"net/http"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrtype"
)

const (
	maxValuePageSize = 1000000
)

// ParsePageParams  - comment func.
func ParsePageParams(r *http.Request, keyIndex, keySize string) (mrtype.PageParams, error) {
	index, err := ParseInt64(r, keyIndex, false)
	if err != nil {
		return mrtype.PageParams{}, err
	}

	size, err := ParseInt64(r, keySize, false)
	if err != nil {
		return mrtype.PageParams{}, err
	}

	if index < 0 {
		index = 0
	}

	if size < 0 {
		size = 0
	}

	if size > maxValuePageSize {
		return mrtype.PageParams{}, mrcore.ErrHttpRequestParamMax.New(keySize, maxValuePageSize)
	}

	return mrtype.PageParams{
		Index: uint64(index),
		Size:  uint64(size),
	}, nil
}
