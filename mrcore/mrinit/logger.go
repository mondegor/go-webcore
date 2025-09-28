package mrinit

import (
	"context"
	sdkslog "log/slog"
	"os"
	"strings"

	"github.com/mondegor/go-sysmess/mrerrors"
	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-sysmess/mrlog/color"
	"github.com/mondegor/go-sysmess/mrlog/slog"
	"github.com/mondegor/go-sysmess/mrlog/slog/middleware"
	"github.com/mondegor/go-sysmess/mrtrace"

	"github.com/mondegor/go-webcore/mrcore"
)

type (
	// LoggerConfig - comment struct.
	LoggerConfig struct {
		Environment string
		Version     string
		Level       string
		JsonFormat  bool
		TimeFormat  string
		ColorMode   bool
	}

	instantError interface {
		ID() string
		Attrs() []any
	}
)

// InitLogger - создаёт и инициализирует логгер,
// помещает его в контекст и возвращает этот контекст.
func InitLogger(cfg LoggerConfig) (logger mrlog.Logger, err error) {
	logger, err = newLogger(cfg)
	if err != nil {
		return nil, err
	}

	if !strings.HasPrefix(cfg.Environment, "local") {
		logger = logger.WithAttrs(mrcore.KeyAppEnvironment, cfg.Environment)

		if cfg.Version != "" {
			logger = logger.WithAttrs(mrcore.KeyAppVersion, cfg.Version)
		}
	}

	return logger, nil
}

func newLogger(cfg LoggerConfig) (*slog.LoggerAdapter, error) {
	if cfg.Level == "" {
		cfg.Level = "info"
	}

	if cfg.TimeFormat == "" {
		cfg.TimeFormat = "RFC3339"
	}

	opts := []slog.Option{
		slog.WithWriter(os.Stdout),
		slog.WithLevel(strings.ToUpper(cfg.Level)),
		slog.WithJsonFormat(cfg.JsonFormat),
		slog.WithTimeFormat(cfg.TimeFormat),
		slog.WithMiddlewareHandler(
			middleware.BeforeHandle(
				func(ctx context.Context, rec sdkslog.Record) sdkslog.Record {
					rec.Attrs(func(attr sdkslog.Attr) bool {
						if attr.Value.Kind() == sdkslog.KindAny {
							if e, ok := attr.Value.Any().(instantError); ok {
								if id := e.ID(); id != "" {
									rec.Add(mrcore.KeyErrorID, id)
								}

								rec.Add(e.Attrs()...)
							}
						}

						return true
					})

					rec.Add(mrtrace.ExtractKeysValues(ctx)...)

					return rec
				},
			),
		),
		slog.WithReplaceAttrs(func(attr sdkslog.Attr) (newAttr sdkslog.Attr) {
			if attr.Value.Kind() == sdkslog.KindAny {
				if e, ok := attr.Value.Any().(*mrerrors.InstantError); ok {
					attr.Value = sdkslog.AnyValue(mrerrors.CastLessVerboseError(e))
				}
			}

			return attr
		}),
		slog.WithColorMode(cfg.ColorMode),
	}

	if cfg.ColorMode {
		opts = append(
			opts,
			slog.WithColorizeAttr(mrcore.KeyAppEnvironment, color.Yellow, color.LightGray),
			slog.WithColorizeAttr(mrcore.KeyAppVersion, color.Yellow, color.LightGray),
			slog.WithColorizeAttr(mrcore.KeyErrorID, color.Yellow, color.Red),

			slog.WithColorizeAttr(mrtrace.KeyProcessID, color.Yellow, color.LightGray),
			slog.WithColorizeAttr(mrtrace.KeyWorkerID, color.Yellow, color.LightGray),

			slog.WithColorizeAttr(mrtrace.KeyCorrelationID, color.LightYellow, color.LightGray),
			slog.WithColorizeAttr(mrtrace.KeyTaskID, color.LightYellow, color.LightGray),
			slog.WithColorizeAttr(mrtrace.KeyRequestID, color.LightYellow, color.LightGray),

			slog.WithColorizeAttr("sql", color.Cyan, color.Green),
		)
	}

	return slog.NewLoggerAdapter(opts...)
}
