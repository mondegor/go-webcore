package mridempotency

type (
	// Responser - интерфейс возвращения результата провайдером идемпотентности.
	Responser interface {
		Content() []byte
		StatusCode() int
	}
)
