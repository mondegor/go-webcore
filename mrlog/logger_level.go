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
	// Level - уровень логирования.
	Level int8
)

var (
	levelName = map[Level]string{
		DebugLevel: "DEBUG",
		InfoLevel:  "INFO",
		WarnLevel:  "WARN",
		ErrorLevel: "ERROR",
		FatalLevel: "FATAL",
		TraceLevel: "TRACE",
	}

	levelValue = map[string]Level{
		"DEBUG": DebugLevel,
		"INFO":  InfoLevel,
		"WARN":  WarnLevel,
		"ERROR": ErrorLevel,
		"FATAL": FatalLevel,
		"TRACE": TraceLevel,
	}
)

// ParseLevel - comment func.
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
