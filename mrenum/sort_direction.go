package mrenum

import (
	"github.com/mondegor/go-sysmess/mrerr/mr"
)

// Направление сортировки.
const (
	SortDirectionASC  SortDirection = iota // по возрастанию
	SortDirectionDESC                      // по убыванию
)

const (
	enumNameSortDirection = "SortDirection"
)

type (
	// SortDirection - направление сортировки.
	SortDirection uint8
)

var (
	sortDirectionName = map[SortDirection]string{ //nolint:gochecknoglobals
		SortDirectionASC:  "ASC",
		SortDirectionDESC: "DESC",
	}

	sortDirectionValue = map[string]SortDirection{ //nolint:gochecknoglobals
		"ASC":  SortDirectionASC,
		"DESC": SortDirectionDESC,
	}
)

// ParseAndSet - парсит указанное значение и если оно валидно, то устанавливает его числовое значение.
func (e *SortDirection) ParseAndSet(value string) error {
	if parsedValue, ok := sortDirectionValue[value]; ok {
		*e = parsedValue

		return nil
	}

	return mr.ErrInternalKeyNotFoundInSource.New(value, enumNameSortDirection)
}

// String - возвращает значение в виде строки.
func (e SortDirection) String() string {
	return sortDirectionName[e]
}
