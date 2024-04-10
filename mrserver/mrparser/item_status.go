package mrparser

import (
	"net/http"

	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/mrreq"
)

type (
	ItemStatus struct {
		defaultItems []mrenum.ItemStatus
	}
)

var (
	// Make sure the ItemStatus conforms with the mrserver.RequestParserItemStatus interface
	_ mrserver.RequestParserItemStatus = (*ItemStatus)(nil)
)

func NewItemStatus() *ItemStatus {
	return &ItemStatus{}
}

func NewItemStatusWithDefault(items []mrenum.ItemStatus) *ItemStatus {
	return &ItemStatus{
		defaultItems: items,
	}
}

func (p *ItemStatus) FilterStatusList(r *http.Request, key string) []mrenum.ItemStatus {
	items, err := p.parseList(r, key)

	if err != nil {
		mrlog.Ctx(r.Context()).Warn().Err(err).Send()
		return p.defaultItems
	}

	if len(items) == 0 {
		return p.defaultItems
	}

	return items
}

func (p *ItemStatus) parseList(r *http.Request, key string) ([]mrenum.ItemStatus, error) {
	enumList, err := mrreq.ParseEnumList(r, key)

	if err != nil {
		return []mrenum.ItemStatus{}, err
	}

	items, err := mrenum.ParseItemStatusList(enumList)

	if err != nil {
		return []mrenum.ItemStatus{}, err
	}

	return items, nil
}
