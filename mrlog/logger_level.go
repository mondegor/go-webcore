package mrlog

import (
	"github.com/mondegor/go-webcore/mrcore"
)

const (
	DebugLevel Level = iota // DebugLevel - FatalLevel + WarnLevel + WarnLevel + InfoLevel + отладочные сообщения
	InfoLevel               // InfoLevel - FatalLevel + WarnLevel + WarnLevel + информационные сообщения
	WarnLevel               // WarnLevel - FatalLevel + WarnLevel + предупреждения
	ErrorLevel              // ErrorLevel - FatalLevel + ошибки
	FatalLevel              // FatalLevel - отображение только критических ошибок
	TraceLevel Level = -1   // TraceLevel - FatalLevel + WarnLevel + WarnLevel + InfoLevel + DebugLevel + трассировочные сообщения
)

type (
	// Level - comment type.
	Level int8
)

var (
	levelName = map[Level]string{
		DebugLevel: "debug",
		InfoLevel:  "info",
		WarnLevel:  "warn",
		ErrorLevel: "error",
		FatalLevel: "fatal",
		TraceLevel: "trace",
	}

	levelValue = map[string]Level{
		"debug": DebugLevel,
		"info":  InfoLevel,
		"warn":  WarnLevel,
		"error": ErrorLevel,
		"fatal": FatalLevel,
		"trace": TraceLevel,
	}
)

// ParseLevel  - comment func.
func ParseLevel(str string) (Level, error) {
	if value, ok := levelValue[str]; ok {
		return value, nil
	}

	return InfoLevel, mrcore.ErrInternalKeyNotFoundInSource.New(str, "mrlog.Level")
}

// String - comment method.
func (e Level) String() string {
	return levelName[e]
}
