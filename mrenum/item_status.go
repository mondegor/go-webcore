package mrenum

import (
	"database/sql/driver"
	"encoding/json"
	"math"

	"github.com/mondegor/go-sysmess/mrerr/mr"
)

// Статусы элемента.
const (
	ItemStatusDraft    ItemStatus = iota + 1 // черновик
	ItemStatusEnabled                        // действующий
	ItemStatusDisabled                       // отключённый
)

const (
	itemStatusLast     = uint8(ItemStatusDisabled)
	enumNameItemStatus = "ItemStatus"
)

type (
	// ItemStatus - статус элемента.
	ItemStatus uint8
)

var (
	itemStatusName = map[ItemStatus]string{ //nolint:gochecknoglobals
		ItemStatusDraft:    "DRAFT",
		ItemStatusEnabled:  "ENABLED",
		ItemStatusDisabled: "DISABLED",
	}

	itemStatusValue = map[string]ItemStatus{ //nolint:gochecknoglobals
		"DRAFT":    ItemStatusDraft,
		"ENABLED":  ItemStatusEnabled,
		"DISABLED": ItemStatusDisabled,
	}
)

// ParseAndSet - парсит указанное значение и если оно валидно, то устанавливает его числовое значение.
func (e *ItemStatus) ParseAndSet(value string) error {
	if parsedValue, ok := itemStatusValue[value]; ok {
		*e = parsedValue

		return nil
	}

	return mr.ErrInternalKeyNotFoundInSource.New(value, enumNameItemStatus)
}

// Set - устанавливает указанное значение, если оно является enum значением.
func (e *ItemStatus) Set(value uint8) error {
	if value > 0 && value <= itemStatusLast {
		*e = ItemStatus(value)

		return nil
	}

	return mr.ErrInternalKeyNotFoundInSource.New(value, enumNameItemStatus)
}

// String - возвращает значение в виде строки.
func (e ItemStatus) String() string {
	return itemStatusName[e]
}

// // Empty - сообщает, установлено ли enum значение.
// func (e ItemStatus) Empty() bool {
// 	return e == 0
// }

// MarshalJSON - переводит enum значение в строковое представление.
func (e ItemStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.String())
}

// UnmarshalJSON - переводит строковое значение в enum представление.
func (e *ItemStatus) UnmarshalJSON(data []byte) error {
	var value string

	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	return e.ParseAndSet(value)
}

// Scan implements the Scanner interface.
func (e *ItemStatus) Scan(value any) error {
	if val, ok := value.(int64); ok && val >= 0 && val <= math.MaxUint8 {
		return e.Set(uint8(val)) //nolint:gosec
	}

	return mr.ErrInternalTypeAssertion.New(enumNameItemStatus, value)
}

// Value implements the driver.Valuer interface.
func (e ItemStatus) Value() (driver.Value, error) {
	return uint8(e), nil
}

// ParseItemStatusList - парсит массив строковых значений и
// возвращает соответствующий массив enum значений.
func ParseItemStatusList(items []string) ([]ItemStatus, error) {
	var tmp ItemStatus

	parsedItems := make([]ItemStatus, len(items))

	for i := range items {
		if err := tmp.ParseAndSet(items[i]); err != nil {
			return nil, err
		}

		parsedItems[i] = tmp
	}

	return parsedItems, nil
}
