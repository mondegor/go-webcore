package mrreq

import (
	"net/http"

	"github.com/mondegor/go-webcore/mrenum"
)

func ParseItemStatusList(r *http.Request, key string, defaultItem mrenum.ItemStatus) ([]mrenum.ItemStatus, error) {
	def := func (defaultItem mrenum.ItemStatus) []mrenum.ItemStatus {
		if defaultItem == 0 {
			return []mrenum.ItemStatus{}
		}

		return []mrenum.ItemStatus{defaultItem}
	}

	enums, err := ParseEnumList(r, key)

	if err != nil {
		return def(defaultItem), err
	}

	items, err := mrenum.ParseItemStatusList(enums)

	if err != nil {
		return def(defaultItem), err
	}

	if len(items) == 0 {
		return def(defaultItem), nil
	}

	return items, nil
}
