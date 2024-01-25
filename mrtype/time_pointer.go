package mrtype

import "time"

func TimePointer(value time.Time) *time.Time {
	if value.IsZero() {
		return nil
	}

	return &value
}

func TimePointerCopy(value *time.Time) *time.Time {
	if value == nil || value.IsZero() {
		return nil
	}

	c := *value

	return &c
}
