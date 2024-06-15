package mrlogbase

import (
	"fmt"
)

func appendKey(dst []byte, key string) []byte {
	if len(dst) == 0 {
		dst = append(dst, ' ')
	} else {
		dst = append(dst, ',', ' ')
	}

	if key == "" {
		key = "key"
	}

	return append(append(dst, key...), '=')
}

func appendValue(dst []byte, value any) []byte {
	return append(dst, fmt.Sprintf("\"%v\"", value)...)
}
