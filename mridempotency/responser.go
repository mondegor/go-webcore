package mridempotency

type (
	// Responser - определяет интерфейс для возврата сохранённого ответа
	// из хранилища идемпотентности.
	Responser interface {
		Content() []byte
		StatusCode() int
	}
)
