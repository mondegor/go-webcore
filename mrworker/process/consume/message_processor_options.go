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

// WithReadyTimeout - устанавливает опцию readyTimeout для MessageProcessor.
func WithReadyTimeout(value time.Duration) Option {
	return func(p *MessageProcessor) {
		if value > 0 {
			p.readyTimeout = value
		}
	}
}

// WithStartReadDelay - устанавливает опцию отложенного старта чтения данных.
func WithStartReadDelay(value time.Duration) Option {
	return func(p *MessageProcessor) {
		p.startReadDelay = value
	}
}

// WithReadPeriod - устанавливает опцию периода чтения данных консьюмером, когда он в состоянии простоя.
func WithReadPeriod(value time.Duration) Option {
	return func(p *MessageProcessor) {
		if value > 0 {
			p.readPeriod = value
		}
	}
}

// WithCancelReadTimeout - устанавливает опцию таймаута на время отмены чтения данных
// консьюмером при неожиданном завершении работы воркеров.
func WithCancelReadTimeout(value time.Duration) Option {
	return func(p *MessageProcessor) {
		if value > 0 {
			p.cancelReadTimeout = value
		}
	}
}

// WithHandlerTimeout - устанавливает опцию handlerTimeout выполнения обработчика сообщения.
func WithHandlerTimeout(value time.Duration) Option {
	return func(p *MessageProcessor) {
		if value > 0 {
			p.handlerTimeout = value
		}
	}
}

// WithQueueSize - устанавливает опцию размера очереди обработки сообщений.
func WithQueueSize(value uint32) Option {
	return func(p *MessageProcessor) {
		if p.queueSize > 0 {
			p.queueSize = value
		}
	}
}

// WithWorkersCount - устанавливает опцию количества воркеров обрабатывающих сообщения.
func WithWorkersCount(value uint16) Option {
	return func(p *MessageProcessor) {
		if value > 0 {
			p.workersCount = value
		}
	}
}
