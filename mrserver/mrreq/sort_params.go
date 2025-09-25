package mrreq

import (
	"regexp"
	"strings"

	"github.com/mondegor/go-sysmess/mrerr/mr"

	"github.com/mondegor/go-webcore/mrtype"
)

const (
	maxLenSortField = 32
)

var regexpSorterField = regexp.MustCompile(`^[a-z]([a-zA-Z0-9]+)?[a-zA-Z0-9]$`)

// ParseSortParams - возвращает SortParams из строковых параметров по указанным ключам.
func ParseSortParams(getter valueGetter, keyField, keyDirection string) (mrtype.SortParams, error) {
	value := getter.Get(keyField)

	if value == "" {
		return mrtype.SortParams{}, nil
	}

	if len(value) > maxLenSortField {
		return mrtype.SortParams{}, mr.ErrHttpRequestParamLenMax.New(keyField, maxLenSortField)
	}

	if !regexpSorterField.MatchString(value) {
		return mrtype.SortParams{}, mr.ErrHttpRequestParseParam.New(keyField, "SortParams", value)
	}

	var params mrtype.SortParams

	if direction := getter.Get(keyDirection); direction != "" {
		if err := params.Direction.ParseAndSet(strings.ToUpper(direction)); err != nil {
			return params, err
		}
	}

	params.FieldName = value

	return params, nil
}
