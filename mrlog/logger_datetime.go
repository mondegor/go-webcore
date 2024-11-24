package mrlog

import (
	"fmt"
	"time"

	"github.com/mondegor/go-webcore/mrcore"
)

// ParseDateTimeFormat - comment func.
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

	return time.RFC3339, mrcore.ErrInternalWithDetails.New(fmt.Sprintf("value '%s' is not found in mrlog.ParseDateTimeFormat()", str))
}
