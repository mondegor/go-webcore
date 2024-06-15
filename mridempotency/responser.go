package mridempotency

type (
	// Responser - интерфейс возвращения результата провайдером идемпотентности.
	Responser interface {
		StatusCode() int
		Body() []byte
	}
)
