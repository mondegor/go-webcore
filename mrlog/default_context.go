package mrlog

type (
	defaultContext struct {
		logger DefaultLogger
	}
)

func (c *defaultContext) Logger() Logger {
	return &c.logger
}

func (c *defaultContext) CallerWithSkipFrame(count int) LoggerContext {
	cp := *c
	return &cp
}

func (c *defaultContext) Str(key, value string) LoggerContext {
	cp := *c
	cp.logger.context = appendValue(appendKey(cp.logger.context, key), value)
	return &cp
}

func (c *defaultContext) Bytes(key string, value []byte) LoggerContext {
	cp := *c
	cp.logger.context = appendValue(appendKey(cp.logger.context, key), value)
	return &cp
}

func (c *defaultContext) Int(key string, value int) LoggerContext {
	cp := *c
	cp.logger.context = appendValue(appendKey(cp.logger.context, key), value)
	return &cp
}

func (c *defaultContext) Any(key string, value any) LoggerContext {
	cp := *c
	cp.logger.context = appendValue(appendKey(cp.logger.context, key), value)
	return &cp
}
