package nopresponser

import "net/http"

type (
	// Responser - заглушка реализующая интерфейс возвращения
	// результата провайдером идемпотентности.
	Responser struct{}
)

// New - создаёт объект Responser.
func New() *Responser {
	return &Responser{}
}

// StatusCode - возвращает код ответа 200.
func (r Responser) StatusCode() int {
	return http.StatusOK
}

// Content - возвращает пустое содержание.
func (r Responser) Content() []byte {
	return nil
}
