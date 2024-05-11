package mrlog

import (
	"fmt"
	"time"
)

func ParseDateTimeFormat(str string) (string, error) {
	switch str {
	case "RFC3339":
		return time.RFC3339, nil
	case "RFC3339Nano":
		return time.RFC3339Nano, nil
	case "DateTime":
		return time.DateTime, nil
	case "TimeOnly":
		return time.TimeOnly, nil

	}

	return time.RFC3339, fmt.Errorf("'%s' is not found in mrlog.ParseDateTimeFormat()", str)
}
