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

		// TODO: не безопасно!!! можно использовать только
		// TODO: при инициализации в глобальном потоке и только один раз!!!
		if zerolog.TimeFieldFormat != opts.TimestampFormat {
			zerolog.TimeFieldFormat = opts.TimestampFormat
		}
	}

	var out io.Writer = os.Stdout

	if !opts.JsonFormat {
		if opts.TimestampFormat == "" {
			opts.TimestampFormat = time.Kitchen
		}

		out = zerolog.ConsoleWriter{
			Out:        out,
			TimeFormat: opts.TimestampFormat,
			NoColor:    !opts.ConsoleColor,
		}
	}

	zero := zerolog.New(out)

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
