package mrinit

import "github.com/mondegor/go-sysmess/mrerr"

type (
	// ErrorSettings - настройки ошибки, которые должны быть применены по совпадению кода ошибки.
	ErrorSettings struct {
		Code          string
		WithCaller    bool
		WithOnCreated bool
	}

	codeGetter interface {
		Code() string
	}
)

// AllEnabled - возвращает настройки, при которых всегда формируется
// стек вызовов и генерируется событие создания ошибки.
func AllEnabled(value codeGetter) ErrorSettings {
	return ErrorSettings{
		Code:          value.Code(),
		WithCaller:    true,
		WithOnCreated: true,
	}
}

// WithCaller - возвращает настройки, при которых всегда формируется
// стек вызовов и отключается генерация события создания ошибки.
func WithCaller(value codeGetter) ErrorSettings {
	return ErrorSettings{
		Code:          value.Code(),
		WithCaller:    false,
		WithOnCreated: true,
	}
}

// WithOnCreated - возвращает настройки, при которых всегда генерируется
// событие создания ошибки и отключается формирование стека вызовов.
func WithOnCreated(value codeGetter) ErrorSettings {
	return ErrorSettings{
		Code:          value.Code(),
		WithCaller:    true,
		WithOnCreated: false,
	}
}

// AllDisabled - возвращает настройки, при которых все опции отключены.
func AllDisabled(value codeGetter) ErrorSettings {
	return ErrorSettings{
		Code:          value.Code(),
		WithCaller:    false,
		WithOnCreated: false,
	}
}

// CreateErrorOptionsMap - формирует сопоставление кода ошибки и опций по умолчанию,
// которые этой ошибке должны быть присвоены.
func CreateErrorOptionsMap(
	settings []ErrorSettings,
	callerOption mrerr.ProtoOption,
	onCreatedOption mrerr.ProtoOption,
) (code2options map[string][]mrerr.ProtoOption) {
	if len(settings) == 0 {
		return nil
	}

	code2options = make(map[string][]mrerr.ProtoOption, len(settings))

	for _, es := range settings {
		var options []mrerr.ProtoOption

		if es.WithCaller && callerOption != nil {
			options = append(options, callerOption)
		}

		if es.WithOnCreated && onCreatedOption != nil {
			options = append(options, onCreatedOption)
		}

		code2options[es.Code] = options
	}

	return code2options
}
