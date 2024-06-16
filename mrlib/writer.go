package mrlib

import (
	"context"
	"io"

	"github.com/mondegor/go-webcore/mrlog"
)

// Write - comment func.
func Write(ctx context.Context, w io.Writer, bytes []byte) {
	if len(bytes) < 1 {
		return
	}

	if _, err := w.Write(bytes); err != nil {
		mrlog.Ctx(ctx).Error().Err(err).Msg("write failed")
	}
}
