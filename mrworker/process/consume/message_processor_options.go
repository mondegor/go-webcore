package consume

import (
	"time"
)

type (
	// Option - настройка объекта MessageProcessor.
	Option func(p *MessageProcessor)
)

// WithCaption - устанавливает опцию caption для MessageProcessor.
func WithCaption(value string) Option {
	return func(p *MessageProcessor) {
		if value != "" {
			p.caption = value
		}
	}
}

// WithPeriod - устанавливает опцию периода обращения к консьюмеру для MessageProcessor.
func WithPeriod(value time.Duration) Option {
	return func(p *MessageProcessor) {
		if value > 0 {
			p.period = value
		}
	}
}

// WithTimeout - устанавливает опцию timeout выполнения
// обработчика сообщения для MessageProcessor.
func WithTimeout(value time.Duration) Option {
	return func(p *MessageProcessor) {
		if value > 0 {
			p.timeout = value
		}
	}
}

// WithQueueSize - устанавливает опцию размера очереди обработки сообщений для MessageProcessor.
func WithQueueSize(value uint32) Option {
	return func(p *MessageProcessor) {
		p.queueSize = value
	}
}

// WithWorkersCount - устанавливает опцию количества воркеров обрабатывающих сообщения для MessageProcessor.
func WithWorkersCount(value uint16) Option {
	return func(p *MessageProcessor) {
		if value > 0 {
			p.workersCount = value
		}
	}
}
