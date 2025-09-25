package mrlib

import (
	"time"
)

// TimeLeftInSec - возвращает количество оставшихся секунд
// до наступления указанного времени (округляется в большую сторону).
// Если результат получился отрицательным, возвращается 0.
func TimeLeftInSec(tm time.Time) int64 {
	const (
		second     = 1e9
		halfSecond = second / 2
	)

	rest := tm.UnixNano() - time.Now().UnixNano() + halfSecond

	if rest < 0 {
		return 0
	}

	return rest / second
}
