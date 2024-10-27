package schedule

import "github.com/mondegor/go-webcore/mrworker"

type (
	// Option - настройка объекта TaskScheduler.
	Option func(s *TaskScheduler)
)

// WithCaption - устанавливает опцию caption для TaskScheduler.
func WithCaption(value string) Option {
	return func(s *TaskScheduler) {
		if value != "" {
			s.caption = value
		}
	}
}

// WithTasks - устанавливает опцию tasks для TaskScheduler.
func WithTasks(values ...mrworker.Task) Option {
	return func(s *TaskScheduler) {
		s.tasks = values
	}
}
