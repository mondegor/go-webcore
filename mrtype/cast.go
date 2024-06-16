package mrtype

import (
	"time"
)

// BoolToInt64 - преобразование bool к int64.
func BoolToInt64(value bool) int64 {
	if value {
		return 1
	}

	return 0
}

// BoolToPointer - required = false by default.
func BoolToPointer(value bool, required ...bool) *bool {
	if isNullable(required) && !value {
		return nil
	}

	return &value
}

// Int32ToPointer - required = false by default.
func Int32ToPointer(value int32, required ...bool) *int32 {
	if isNullable(required) && value == 0 {
		return nil
	}

	return &value
}

// Int64ToPointer - required = false by default.
func Int64ToPointer(value int64, required ...bool) *int64 {
	if isNullable(required) && value == 0 {
		return nil
	}

	return &value
}

// StringToPointer - required = false by default.
func StringToPointer(value string, required ...bool) *string {
	if isNullable(required) && value == "" {
		return nil
	}

	return &value
}

// TimeToPointer - required = false by default.
func TimeToPointer(value time.Time, required ...bool) *time.Time {
	if isNullable(required) && value.IsZero() {
		return nil
	}

	return &value
}

// TimePointerCopy - comment func.
func TimePointerCopy(value *time.Time) *time.Time {
	if value == nil || value.IsZero() {
		return nil
	}

	c := *value

	return &c
}

func isNullable(required []bool) bool {
	return len(required) < 1 || !required[0]
}
