package mrparser

import (
	"net/http"

	"github.com/mondegor/go-sysmess/mrlog"

	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrserver/mrreq"
)

type (
	// ItemStatus - парсер mrenum.ItemStatus.
	ItemStatus struct {
		logger       mrlog.Logger
		defaultItems []mrenum.ItemStatus
	}
)

// NewItemStatus - создаёт объект ItemStatus.
func NewItemStatus(logger mrlog.Logger) *ItemStatus {
	return &ItemStatus{
		logger: logger,
	}
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
		p.logger.Warn(r.Context(), "FilterStatusList", "error", err)

		return p.defaultItems
	}

	if len(items) == 0 {
		return p.defaultItems
	}

	return items
}

func (p *ItemStatus) parseList(r *http.Request, key string) ([]mrenum.ItemStatus, error) {
	enumList, err := mrreq.ParseEnumList(r.URL.Query(), key)
	if err != nil {
		return nil, err
	}

	items, err := mrenum.ParseItemStatusList(enumList)
	if err != nil {
		return nil, err
	}

	return items, nil
}
