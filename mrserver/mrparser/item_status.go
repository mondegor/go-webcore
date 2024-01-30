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
	}
)

// Make sure the ItemStatus conforms with the mrserver.RequestParserItemStatus interface
var _ mrserver.RequestParserItemStatus = (*ItemStatus)(nil)

func NewItemStatus() *ItemStatus {
	return &ItemStatus{}
}

func (p *ItemStatus) FilterStatusList(r *http.Request, key string) []mrenum.ItemStatus {
	items, err := mrreq.ParseItemStatusList(
		r,
		key,
		[]mrenum.ItemStatus{
			// :TODO: добавить значение по умолчнию в настройки
			mrenum.ItemStatusEnabled,
		},
	)

	if err != nil {
		mrlog.Ctx(r.Context()).Warn().Err(err)
		return []mrenum.ItemStatus{}
	}

	return items
}
