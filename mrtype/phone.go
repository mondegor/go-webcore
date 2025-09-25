package mrtype

import (
	"strconv"
)

// UintToPhone - comment func.
func UintToPhone(value uint64) string {
	if value > 0 {
		return "+" + strconv.FormatUint(value, 10)
	}

	return ""
}
