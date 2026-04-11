package mrserver

const (
	// HeaderKeyAcceptLanguage - заголовок с предпочтительными языками клиента.
	// Пример: "ru-RU,ru;q=0.9,en;q=0.8"
	// Используется для определения языка ответа (переопределяется из профиля пользователя).
	HeaderKeyAcceptLanguage = "Accept-Language"

	// HeaderKeyCorrelationID - внешний ID цепочки запросов для распределённой трассировки.
	// Передаётся клиентом для связывания запросов через несколько сервисов.
	// Пример: f7479171-83d2-4f64-84ac-892f8c0aaf48.
	HeaderKeyCorrelationID = "X-Correlation-Id"

	// HeaderKeyIdempotencyKey - ключ идемпотентности для предотвращения повторных операций.
	// Пример: 12b779ee-6bfd-495a-94ba-e4fa517f0268.
	HeaderKeyIdempotencyKey = "X-Idempotency-Key"

	// HeaderKeyRequestID - внутренний ID текущего запроса, генерируемый сервером.
	// Используется для трассировки и поиска запроса в логах.
	// Пример: 3c0b88e8-ba21-4d49-afda-d92e7ac08918.
	HeaderKeyRequestID = "X-Request-Id"

	// HeaderKeyUserIDSlashGroup - внутренний заголовок с ID пользователя и его группой.
	// Формат: "{userId}/{realm}/{kind}".
	// Пример: 550e8400-e29b-41d4-a716-446655440000/app/admin.
	HeaderKeyUserIDSlashGroup = "X-Internal-UserId-Group"
)
