package mrparser

import (
	"net/http"

	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/mrreq"
)

type (
	// ItemStatus - парсер mrenum.ItemStatus.
	ItemStatus struct {
		defaultItems []mrenum.ItemStatus
	}
)

// Make sure the ItemStatus conforms with the mrserver.RequestParserItemStatus interface.
var _ mrserver.RequestParserItemStatus = (*ItemStatus)(nil)

// NewItemStatus - создаёт объект ItemStatus.
func NewItemStatus() *ItemStatus {
	return &ItemStatus{}
}

// NewItemStatusWithDefault - создаёт объект ItemStatus со статусами по умолчанию.
func NewItemStatusWithDefault(items []mrenum.ItemStatus) *ItemStatus {
	return &ItemStatus{
		defaultItems: items,
	}
}

// FilterStatusList - возвращает массив mrenum.ItemStatus поступивший из внешнего запроса.
// Если ключ key не найден или возникнет ошибка, то возвращается пустой массив.
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
		return nil, err
	}

	items, err := mrenum.ParseItemStatusList(enumList)
	if err != nil {
		return nil, err
	}

	return items, nil
}
