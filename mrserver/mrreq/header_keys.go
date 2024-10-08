package mrreq

const (
	// HeaderKeyAcceptLanguage - название заголовка содержащий предпочитаемые клиентом языки общения.
	HeaderKeyAcceptLanguage = "Accept-Language"

	// HeaderKeyCorrelationID - название заголовка содержащий пользовательский ID используемый в запросе.
	// sample: f7479171-83d2-4f64-84ac-892f8c0aaf48.
	HeaderKeyCorrelationID = "X-Correlation-Id"

	// HeaderKeyIdempotencyKey - название заголовка содержащего ключ идемпотентности операции.
	HeaderKeyIdempotencyKey = "X-Idempotency-Key"

	// HeaderKeyRequestID - название заголовка содержащего ID текущего запроса.
	// sample: f7479171-83d2-4f64-84ac-892f8c0aaf48, cs0n2utf3kkujsmnq9og.
	HeaderKeyRequestID = "X-Request-Id"
)
