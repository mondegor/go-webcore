package mrrscors

import (
	"net/http"

	"github.com/mondegor/go-webcore/mrlog"
	"github.com/rs/cors"
)

// go get -u github.com/rs/cors

type (
	CorsAdapter struct {
		cors *cors.Cors
	}

	Options struct {
		AllowedOrigins   []string
		AllowedMethods   []string
		AllowedHeaders   []string
		ExposedHeaders   []string
		AllowCredentials bool
		Logger           mrlog.Logger
	}
)

func New(opts Options) *CorsAdapter {
	options := cors.Options{
		AllowedOrigins:   opts.AllowedOrigins,
		AllowedMethods:   opts.AllowedMethods,
		AllowedHeaders:   opts.ExposedHeaders,
		ExposedHeaders:   opts.ExposedHeaders,
		AllowCredentials: opts.AllowCredentials,
	}

	if opts.Logger != nil && opts.Logger.Level() <= mrlog.DebugLevel {
		options.Debug = true

		opts.Logger.Debug().MsgFunc(
			func() string {
				var buf []byte

				buf = append(buf, "Cors.AllowedOrigins:"...)

				for i := range opts.AllowedOrigins {
					buf = append(buf, "\n- "+opts.AllowedOrigins[i]+";"...)
				}

				return string(buf)
			},
		)
	}

	return &CorsAdapter{
		cors: cors.New(options),
	}
}

func (c *CorsAdapter) Middleware(next http.Handler) http.Handler {
	return c.cors.Handler(next)
}
