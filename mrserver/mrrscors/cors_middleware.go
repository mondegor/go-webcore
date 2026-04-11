package mrrscors

import (
	"context"
	"net/http"

	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/rs/cors"
)

// go get -u github.com/rs/cors

type (
	// Options - конфигурация CORS middleware.
	Options struct {
		// AllowedOrigins: список разрешённых origin (например: ["https://example.com", "http://localhost:3000"]).
		AllowedOrigins []string

		// AllowedMethods: список разрешённых HTTP-методов (например: ["GET", "POST", "PUT", "DELETE"]).
		AllowedMethods []string

		// AllowedHeaders: список разрешённых заголовков запроса (например: ["Content-Type", "Authorization"]).
		AllowedHeaders []string

		// ExposedHeaders: список заголовков ответа, доступных клиентскому JS (например: ["X-Request-Id"]).
		ExposedHeaders []string

		// AllowCredentials: разрешить передачу cookies и авторизационных заголовков.
		AllowCredentials bool

		// Logger: логгер для отладки CORS конфигурации.
		Logger mrlog.Logger
	}
)

// Middleware - создаёт HTTP-middleware для обработки CORS-запросов.
// Оборачивает библиотеку github.com/rs/cors в стандартный middleware-интерфейс.
// Автоматически обрабатывает preflight-запросы (OPTIONS) и добавляет CORS-заголовки.
// Возвращает функцию-обёртку для регистрации в роутере.
func Middleware(opts Options) func(next http.Handler) http.Handler {
	options := cors.Options{
		AllowedOrigins:   opts.AllowedOrigins,
		AllowedMethods:   opts.AllowedMethods,
		AllowedHeaders:   opts.AllowedHeaders,
		ExposedHeaders:   opts.ExposedHeaders,
		AllowCredentials: opts.AllowCredentials,
	}

	if opts.Logger == nil {
		opts.Logger = mrlog.NopLogger()
	}

	if mrlog.DebugEnabled(opts.Logger) {
		options.Debug = true

		opts.Logger.DebugFunc(
			context.Background(),
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

	c := cors.New(options)

	return func(next http.Handler) http.Handler {
		return c.Handler(next)
	}
}
