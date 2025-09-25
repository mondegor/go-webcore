package mrlib

import (
	"context"
	"io"

	"github.com/mondegor/go-sysmess/mrlog"
)

// Write - адаптер для вызова io.Writer указанного объекта с логированием ошибки.
func Write(ctx context.Context, logger mrlog.Logger, w io.Writer, bytes []byte) {
	if len(bytes) == 0 {
		return
	}

	if _, err := w.Write(bytes); err != nil {
		logger.Error(ctx, "write failed", "error", err)
	}
}
