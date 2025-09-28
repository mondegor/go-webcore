package mrreq

const (
	// HeaderKeyAcceptLanguage - название заголовка содержащий предпочитаемые клиентом языки общения.
	HeaderKeyAcceptLanguage = "Accept-Language"

	// HeaderKeyCorrelationID - название заголовка содержащий внешний клиентский ID под которым делается запрос.
	// sample: f7479171-83d2-4f64-84ac-892f8c0aaf48 | csoruadf3kkopl6lok80.
	HeaderKeyCorrelationID = "X-Correlation-Id"

	// HeaderKeyIdempotencyKey - название заголовка содержащего ключ идемпотентности операции.
	// sample: 12b779ee-6bfd-495a-94ba-e4fa517f0268
	HeaderKeyIdempotencyKey = "X-Idempotency-Key"

	// HeaderKeyRequestID - название заголовка содержащего ID текущего запроса.
	// sample: 3c0b88e8-ba21-4d49-afda-d92e7ac08918, cs0n2utf3kkujsmnq9og.
	HeaderKeyRequestID = "X-Request-Id"

	// HeaderKeyUserIDSlashGroup - название внутреннего заголовка содержащий ID пользователя и его группу (userId/realm/kind).
	HeaderKeyUserIDSlashGroup = "X-Internal-UserId-Group"
)
