package mrsentry

import (
	"github.com/getsentry/sentry-go"
)

// filterStackTrace - функция срезает верхнюю часть стека вызовов,
// которая не несёт в себе информативности. В массиве bounds указываются
// все названия функций и стек будет срезан по самой нижней из них.
func filterStackTrace(bounds []string) func(event *sentry.Event, _ *sentry.EventHint) *sentry.Event {
	boundMap := make(map[string]bool, len(bounds))
	for _, item := range bounds {
		boundMap[item] = true
	}

	return func(event *sentry.Event, _ *sentry.EventHint) *sentry.Event {
		for _, ex := range event.Exception {
			if ex.Stacktrace == nil {
				continue
			}

			ex.Stacktrace.Frames = filterStackTraceTrimUpper(boundMap, ex.Stacktrace.Frames)
		}

		for _, th := range event.Threads {
			if th.Stacktrace == nil {
				continue
			}

			th.Stacktrace.Frames = filterStackTraceTrimUpper(boundMap, th.Stacktrace.Frames)
		}

		return event
	}
}

func filterStackTraceTrimUpper(boundMap map[string]bool, frames []sentry.Frame) []sentry.Frame {
	for i := 0; i < len(frames); i++ {
		item := frames[i].Module + "." + frames[i].Function

		if _, ok := boundMap[item]; ok {
			if i > 0 {
				return frames[:i]
			}

			return nil
		}
	}

	return frames
}
