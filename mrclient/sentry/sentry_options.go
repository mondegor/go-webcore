package sentry

import (
	"time"

	"github.com/getsentry/sentry-go"
)

type (
	// Option - настройка объекта Adapter.
	Option func(o *options)

	options struct {
		sentryOpts   sentry.ClientOptions
		flushTimeout time.Duration
	}
)

// WithEnvironment - устанавливает опцию Environment для Adapter.
func WithEnvironment(value string) Option {
	return func(o *options) {
		o.sentryOpts.Environment = value
	}
}

// WithRelease - устанавливает опцию Release для Adapter.
func WithRelease(value string) Option {
	return func(o *options) {
		o.sentryOpts.Release = value
	}
}

// WithDebugMode - устанавливает опцию Debug для Adapter.
func WithDebugMode(value bool) Option {
	return func(o *options) {
		o.sentryOpts.Debug = value
	}
}

// WithTracesSampleRate - устанавливает опцию TracesSampleRate для Adapter.
func WithTracesSampleRate(value float64) Option {
	return func(o *options) {
		o.sentryOpts.TracesSampleRate = value
	}
}

// WithFlushTimeout - устанавливает опцию flushTimeout для Adapter.
func WithFlushTimeout(value time.Duration) Option {
	return func(o *options) {
		o.flushTimeout = value
	}
}
