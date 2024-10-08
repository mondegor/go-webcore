package mrlib

import (
	"context"
	"io"

	"github.com/mondegor/go-webcore/mrlog"
)

// Write - адаптер для вызова io.Writer указанного объекта с логированием ошибки.
func Write(ctx context.Context, w io.Writer, bytes []byte) {
	if len(bytes) == 0 {
		return
	}

	if _, err := w.Write(bytes); err != nil {
		mrlog.Ctx(ctx).Error().Err(err).Msg("write failed")
	}
}
