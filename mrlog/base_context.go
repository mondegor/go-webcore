package mrlog

type (
	contextAdapter struct {
		logger BaseLogger
	}
)

// Logger - comment method.
func (c *contextAdapter) Logger() Logger {
	return &c.logger
}

// Str - comment method.
func (c *contextAdapter) Str(key, value string) LoggerContext {
	cp := *c
	cp.logger.context = appendValue(appendKey(cp.logger.context, key), value)

	return &cp
}

// Bytes - comment method.
func (c *contextAdapter) Bytes(key string, value []byte) LoggerContext {
	cp := *c
	cp.logger.context = appendValue(appendKey(cp.logger.context, key), value)

	return &cp
}

// Int - comment method.
func (c *contextAdapter) Int(key string, value int) LoggerContext {
	cp := *c
	cp.logger.context = appendValue(appendKey(cp.logger.context, key), value)

	return &cp
}

// Any - comment method.
func (c *contextAdapter) Any(key string, value any) LoggerContext {
	cp := *c
	cp.logger.context = appendValue(appendKey(cp.logger.context, key), value)

	return &cp
}
