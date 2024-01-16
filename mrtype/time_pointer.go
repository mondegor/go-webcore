package mrtype

import "time"

func TimePointer(value time.Time) *time.Time {
	if value.IsZero() {
		return nil
	}

	return &value
}
