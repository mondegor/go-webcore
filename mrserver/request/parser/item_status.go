package parser

import (
	"net/http"

	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-sysmess/mrstatus/itemstatus"
)

type (
	// ItemStatus - парсер itemstatus.Enum.
	ItemStatus struct {
		*EnumList[itemstatus.Enum]
	}
)

// NewItemStatus - создаёт объект ItemStatus.
func NewItemStatus(logger mrlog.Logger) *ItemStatus {
	return &ItemStatus{
		EnumList: NewEnumList(
			logger,
			itemstatus.ParseList,
		),
	}
}

// NewItemStatusWithDefault - создаёт объект ItemStatus со статусами по умолчанию.
func NewItemStatusWithDefault(logger mrlog.Logger, items []itemstatus.Enum) *ItemStatus {
	return &ItemStatus{
		EnumList: NewEnumListWithDefault(
			logger,
			items,
			itemstatus.ParseList,
		),
	}
}

// FilterStatusList - возвращает массив itemstatus.Enum поступивший из внешнего запроса.
// Если ключ key не найден или возникнет ошибка, то возвращается nil значение.
func (p *ItemStatus) FilterStatusList(r *http.Request, key string) []itemstatus.Enum {
	return p.FilterEnumList(r, key)
}
