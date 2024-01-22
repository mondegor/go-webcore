package mrreq

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrtype"
)

const (
	maxLenSortField = 32
)

var (
	regexpSorterField = regexp.MustCompile(`^[a-z]([a-zA-Z0-9]+)?[a-zA-Z0-9]$`)
)

func ParseSortParams(r *http.Request, keyField, keyDirection string) (mrtype.SortParams, error) {
	query := r.URL.Query()
	value := query.Get(keyField)

	if value == "" {
		return mrtype.SortParams{}, nil
	}

	if len(value) > maxLenSortField {
		return mrtype.SortParams{}, mrcore.FactoryErrHttpRequestParamLenMax.New(keyField, maxLenSortField)
	}

	if !regexpSorterField.MatchString(value) {
		return mrtype.SortParams{}, mrcore.FactoryErrHttpRequestParseParam.New(keyField, "SortParams", value)
	}

	var params mrtype.SortParams

	if direction := query.Get(keyDirection); direction != "" {
		if err := params.Direction.ParseAndSet(strings.ToUpper(direction)); err != nil {
			return params, err
		}
	}

	params.FieldName = value

	return params, nil
}
