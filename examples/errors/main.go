package main

import (
	"os"

	"github.com/mondegor/go-sysmess/mrargs"
	"github.com/mondegor/go-sysmess/mrerr/mr"
	"github.com/mondegor/go-sysmess/mrlog/litelog"
	"github.com/mondegor/go-sysmess/mrlog/slog"
)

func main() {
	l, _ := slog.NewLoggerAdapter(slog.WithWriter(os.Stdout))
	logger := litelog.NewLogger(l)

	logger.Error("this is ErrInternal", "error", mr.ErrInternal)
	logger.Error("this is ErrInternalTypeAssertion", "error", mr.ErrInternalTypeAssertion.New("MY-TYPE", "MY-VALUE"))
	logger.Error("this is ErrInternal with attr", "error", mr.ErrInternal.New().WithAttr("MY-DATA-KEY", mrargs.Group{"itemId": "id-001"}))

	logger.Error(
		"this is ErrInternal with attrs",
		"error", mr.ErrInternal.New(),
		"int1", 1,
		"int2", 2,
		"int3", 3,
	)
}
