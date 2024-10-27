package task

import "time"

type (
	// Option - настройка объекта JobWrapper.
	Option func(j *JobWrapper)
)

// WithCaption - устанавливает опцию caption для JobWrapper.
func WithCaption(value string) Option {
	return func(j *JobWrapper) {
		if value != "" {
			j.caption = value
		}
	}
}

// WithStartup - устанавливает опцию startup для JobWrapper.
func WithStartup(value bool) Option {
	return func(j *JobWrapper) {
		j.startup = value
	}
}

// WithPeriod - устанавливает опцию period для JobWrapper.
func WithPeriod(value time.Duration) Option {
	return func(j *JobWrapper) {
		if value > 0 {
			j.period = value
		}
	}
}

// WithTimeout - устанавливает опцию timeout для JobWrapper.
func WithTimeout(value time.Duration) Option {
	return func(j *JobWrapper) {
		if value > 0 {
			j.timeout = value
		}
	}
}
