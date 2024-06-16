package mrenum

import (
	"github.com/mondegor/go-webcore/mrcore"
)

const (
	SortDirectionASC  SortDirection = iota // SortDirectionASC - сортировка по возрастанию
	SortDirectionDESC                      // SortDirectionDESC - сортировка по убыванию

	enumNameSortDirection = "SortDirection"
)

type (
	// SortDirection - направление сортировки.
	SortDirection uint8
)

var (
	sortDirectionName = map[SortDirection]string{
		SortDirectionASC:  "ASC",
		SortDirectionDESC: "DESC",
	}

	_sortDirectionValue = map[string]SortDirection{
		"ASC":  SortDirectionASC,
		"DESC": SortDirectionDESC,
	}
)

// ParseAndSet - парсит указанное значение и если оно валидно, то устанавливает его числовое значение.
func (e *SortDirection) ParseAndSet(value string) error {
	if parsedValue, ok := _sortDirectionValue[value]; ok {
		*e = parsedValue

		return nil
	}

	return mrcore.ErrInternalKeyNotFoundInSource.New(value, enumNameSortDirection)
}

// String - comment method.
func (e SortDirection) String() string {
	return sortDirectionName[e]
}
