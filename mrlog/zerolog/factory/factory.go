package factory

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"

	"github.com/mondegor/go-webcore/mrlog"
	mrzerolog "github.com/mondegor/go-webcore/mrlog/zerolog"
)

// go get -u github.com/rs/zerolog

type (
	// Options - опции для создания Logger.
	Options struct {
		Stdout           io.Writer
		Level            string
		JsonFormat       bool
		TimestampFormat  string
		ConsoleColor     bool
		PrepareErrorFunc func(err error) error
	}
)

// NewZeroLogAdapter - создаёт объект LoggerAdapter.
func NewZeroLogAdapter(opts Options) (logger *mrzerolog.LoggerAdapter, err error) {
	level, err := mrlog.ParseLevel(opts.Level)
	if err != nil {
		return nil, err
	}

	if opts.TimestampFormat != "" {
		opts.TimestampFormat, err = mrlog.ParseDateTimeFormat(opts.TimestampFormat)
		if err != nil {
			return nil, err
		}

		// TODO: небезопасно!!! можно использовать только
		//       при инициализации в глобальном потоке и только один раз!!!
		if zerolog.TimeFieldFormat != opts.TimestampFormat {
			zerolog.TimeFieldFormat = opts.TimestampFormat
		}
	}

	if opts.Stdout == nil {
		opts.Stdout = os.Stderr
	}

	if !opts.JsonFormat {
		if opts.TimestampFormat == "" {
			opts.TimestampFormat = time.Kitchen
		}

		opts.Stdout = zerolog.ConsoleWriter{
			Out:        opts.Stdout,
			TimeFormat: opts.TimestampFormat,
			NoColor:    !opts.ConsoleColor,
		}
	}

	zero := zerolog.New(opts.Stdout)

	if opts.TimestampFormat != "" {
		zero = zero.With().Timestamp().Logger()
	}

	return mrzerolog.New(
		mrzerolog.Options{
			Zerolog:          zero,
			Level:            level,
			PrepareErrorFunc: opts.PrepareErrorFunc,
		},
	), nil
}
