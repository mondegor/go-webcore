package mrserver

import (
	"net/http"
)

const (
	// VarRestOfURL - шаблонная переменная для обозначения остатка пути в URL.
	// Используется в маршрутах для захвата произвольной части пути после указанного префикса.
	VarRestOfURL = "{{restOfUrl}}"
)

type (
	// HttpRouter - интерфейс маршрутизатора HTTP-запросов.
	//
	// Обеспечивает:
	//  - Регистрацию middleware-обработчиков в цепочку обработки;
	//  - Регистрацию HTTP-контроллеров с их обработчиками;
	//  - Регистрацию отдельных функций-обработчиков для метода и пути;
	//  - Обработку входящих HTTP-запросов (реализация http.Handler);
	HttpRouter interface {
		RegisterMiddleware(handlers ...func(next http.Handler) http.Handler)
		Register(controllers ...HttpController)
		HandlerFunc(method, path string, handler http.HandlerFunc)
		ServeHTTP(w http.ResponseWriter, r *http.Request)
	}

	// HttpController - интерфейс контроллера с набором HTTP-обработчиков.
	// Контроллер - это логическая группа обработчиков, обычно связанных с одной сущностью.
	HttpController interface {
		// Handlers - возвращает список всех обработчиков контроллера.
		Handlers() []HttpHandler
	}

	// HttpHandler - HTTP-обработчик с метаданными маршрутизации и доступа.
	HttpHandler struct {
		// Method - HTTP-метод запроса.
		Method string

		// URL - путь для маршрутизации (может содержать параметры маршрутизатора).
		URL string

		// Permission - разрешение для проверки доступа (пустая строка = без проверки).
		Permission string

		// Func - функция-обработчик запроса.
		Func HttpHandlerFunc
	}

	// HttpHandlerFunc - сигнатура функции HTTP-обработчика.
	// Отличие от стандартного http.HandlerFunc: возвращает ошибку,
	// которую обрабатывает вышестоящий код (например: middleware или адаптер),
	// а не сам обработчик. Это позволяет централизованно обрабатывать ошибки.
	HttpHandlerFunc func(w http.ResponseWriter, r *http.Request) error
)
