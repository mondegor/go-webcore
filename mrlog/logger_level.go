package mrlog

import (
	"fmt"
)

const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel

	TraceLevel Level = -1
)

type (
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

func ParseLevel(str string) (Level, error) {
	if value, ok := levelValue[str]; ok {
		return value, nil
	}

	return InfoLevel, fmt.Errorf("'%s' is not found in map %s", str, "mrlog.Level")
}

func (e Level) String() string {
	return levelName[e]
}
