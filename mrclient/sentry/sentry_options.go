package sentry

import (
	"time"

	"github.com/getsentry/sentry-go"
)

type (
	// Option - настройка объекта Adapter.
	Option func(o *options)

	options struct {
		adapter    *Adapter
		sentryOpts sentry.ClientOptions
	}
)

// WithEnvironment - устанавливает окружение (Environment) для Sentry.
// Используется для разделения ошибок из разных окружений ("production", "testing").
func WithEnvironment(value string) Option {
	return func(o *options) {
		o.sentryOpts.Environment = value
	}
}

// WithRelease - устанавливает версию релиза (Release) для Sentry.
// Позволяет отслеживать ошибки в контексте конкретных версий приложения.
func WithRelease(value string) Option {
	return func(o *options) {
		o.sentryOpts.Release = value
	}
}

// WithDebugMode - включает или отключает режим отладки (Debug) для Sentry.
func WithDebugMode(value bool) Option {
	return func(o *options) {
		o.sentryOpts.Debug = value
	}
}

// WithTracesSampleRate - устанавливает частоту сэмплирования трассировок (TracesSampleRate).
// Значение от 0.0 (нет трассировок) до 1.0 (все трассировки).
func WithTracesSampleRate(value float64) Option {
	return func(o *options) {
		o.sentryOpts.TracesSampleRate = value
	}
}

// WithFlushTimeout - устанавливает таймаут для отправки ожидающих событий при закрытии.
func WithFlushTimeout(value time.Duration) Option {
	return func(o *options) {
		o.adapter.flushTimeout = value
	}
}
