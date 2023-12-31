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
	params := mrtype.SortParams{}
	query := r.URL.Query()

	value := query.Get(keyField)

	if value == "" {
		return params, nil
	}

	if len(value) > maxLenSortField {
		return params, mrcore.FactoryErrHttpRequestParamLenMax.New(keyField, maxLenSortField)
	}

	if !regexpSorterField.MatchString(value) {
		return params, mrcore.FactoryErrHttpRequestParseParam.New(keyField, "SortParams", value)
	}

	direction := query.Get(keyDirection)

	if direction != "" {
		err := params.Direction.ParseAndSet(strings.ToUpper(direction))

		if err != nil {
			return params, err
		}
	}

	params.FieldName = value

	return params, nil
}
