package mrrscors

import (
	"fmt"
	"net/http"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/rs/cors"
)

// go get -u github.com/rs/cors

type (
	CorsAdapter struct {
		cors *cors.Cors
	}

	corsLoggerAdapter struct {
		logger mrcore.Logger
	}

	Options struct {
		AllowedOrigins   []string
		AllowedMethods   []string
		AllowedHeaders   []string
		ExposedHeaders   []string
		AllowCredentials bool
		Logger           mrcore.Logger
	}
)

func New(opt Options) *CorsAdapter {
	options := cors.Options{
		AllowedOrigins:   opt.AllowedOrigins,
		AllowedMethods:   opt.AllowedMethods,
		AllowedHeaders:   opt.ExposedHeaders,
		ExposedHeaders:   opt.ExposedHeaders,
		AllowCredentials: opt.AllowCredentials,
	}

	if opt.Logger != nil && opt.Logger.Level() == mrcore.LogDebugLevel {
		options.Debug = true
		options.Logger = &corsLoggerAdapter{logger: opt.Logger}

		debugInfo := fmt.Sprintf("Cors.AllowedOrigins:")

		for i := range opt.AllowedOrigins {
			debugInfo = fmt.Sprintf(
				"%s\n- %s;",
				debugInfo,
				opt.AllowedOrigins[i],
			)
		}

		opt.Logger.Debug(debugInfo)
	}

	return &CorsAdapter{cors: cors.New(options)}
}

func (c *CorsAdapter) Middleware(next http.Handler) http.Handler {
	return c.cors.Handler(next)
}

func (l *corsLoggerAdapter) Printf(message string, args ...any) {
	l.logger.Debug(message, args...)
}
