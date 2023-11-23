package mrserver

import (
	"fmt"
	"net/http"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/rs/cors"
)

// go get -u github.com/rs/cors

type (
	corsAdapter struct {
		cors *cors.Cors
	}

	CorsOptions struct {
		AllowedOrigins   []string
		AllowedMethods   []string
		AllowedHeaders   []string
		ExposedHeaders   []string
		AllowCredentials bool
		Debug            bool
	}
)

func NewCors(opt CorsOptions) *corsAdapter {
	debugInfo := fmt.Sprintf("Cors.AllowedOrigins:")

	for i := range opt.AllowedOrigins {
		debugInfo = fmt.Sprintf(
			"%s\n- %s;",
			debugInfo,
			opt.AllowedOrigins[i],
		)
	}

	mrcore.LogDebug(debugInfo)

	return &corsAdapter{
		cors: cors.New(cors.Options{
			AllowedOrigins:   opt.AllowedOrigins,
			AllowedMethods:   opt.AllowedMethods,
			AllowedHeaders:   opt.ExposedHeaders,
			ExposedHeaders:   opt.ExposedHeaders,
			AllowCredentials: opt.AllowCredentials,
			Debug:            opt.Debug,
		})}
}

func (c *corsAdapter) Middleware(next http.Handler) http.Handler {
	return c.cors.Handler(next)
}
