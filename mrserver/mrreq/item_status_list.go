package mrreq

import (
	"net/http"

	"github.com/mondegor/go-webcore/mrenum"
)

func ParseItemStatusList(r *http.Request, key string, defaultItems []mrenum.ItemStatus) ([]mrenum.ItemStatus, error) {
	def := func(defaultItems []mrenum.ItemStatus) []mrenum.ItemStatus {
		if len(defaultItems) == 0 {
			return []mrenum.ItemStatus{}
		}

		return defaultItems
	}

	enums, err := ParseEnumList(r, key)

	if err != nil {
		return def(defaultItems), err
	}

	items, err := mrenum.ParseItemStatusList(enums)

	if err != nil {
		return def(defaultItems), err
	}

	if len(items) == 0 {
		return def(defaultItems), nil
	}

	return items, nil
}
