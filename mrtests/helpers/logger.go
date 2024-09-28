package helpers

import (
	"context"

	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrlog/noplog"
)

// ContextWithNopLogger - возвращает контекст с логерром-заглушкой.
func ContextWithNopLogger() context.Context {
	return mrlog.WithContext(context.Background(), noplog.New())
}
