package mrenum

import (
	"fmt"
)

const (
	SortDirectionASC SortDirection = iota
	SortDirectionDESC

	enumNameSortDirection = "SortDirection"
)

type (
	SortDirection uint8
)

var (
	sortDirectionName = map[SortDirection]string{
		SortDirectionASC:  "ASC",
		SortDirectionDESC: "DESC",
	}

	sortDirectionValue = map[string]SortDirection{
		"ASC":  SortDirectionASC,
		"DESC": SortDirectionDESC,
	}
)

func (e *SortDirection) ParseAndSet(value string) error {
	if parsedValue, ok := sortDirectionValue[value]; ok {
		*e = parsedValue
		return nil
	}

	return fmt.Errorf("'%s' is not found in map %s", value, enumNameSortDirection)
}

func (e SortDirection) String() string {
	return sortDirectionName[e]
}
