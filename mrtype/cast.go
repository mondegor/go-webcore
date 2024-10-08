package mrtype

import (
	"time"

	"golang.org/x/exp/constraints"
)

// CastBoolToNumber - возвращает преобразованный bool к Number.
func CastBoolToNumber[Number constraints.Integer | constraints.Float](value bool) Number {
	if value {
		return 1
	}

	return 0
}

// CastBoolToPointer - возвращает преобразованный bool к его указателю.
// Если свойство required не указано или равно false, то вместо значения false будет возвращено nil.
func CastBoolToPointer(value bool, required ...bool) *bool {
	if isNullable(required) && !value {
		return nil
	}

	return &value
}

// CastNumberToPointer - возвращает преобразованный number к его указателю.
// Если свойство required не указано или равно false, то вместо нулевого значения будет возвращено nil.
func CastNumberToPointer[Number constraints.Integer | constraints.Float](value Number, required ...bool) *Number {
	if isNullable(required) && value == 0 {
		return nil
	}

	return &value
}

// CastStringToPointer - возвращает преобразованную строку к его указателю.
// Если свойство required не указано или равно false, то вместо пустого значения будет возвращено nil.
func CastStringToPointer(value string, required ...bool) *string {
	if isNullable(required) && value == "" {
		return nil
	}

	return &value
}

// CastTimeToPointer - возвращает преобразованное время к его указателю.
// Если свойство required не указано или равно false, то вместо нулевого значения будет возвращено nil.
func CastTimeToPointer(value time.Time, required ...bool) *time.Time {
	if isNullable(required) && value.IsZero() {
		return nil
	}

	return &value
}

// CopyTimePointer - возвращает копию значения времени или nil если значение равно nil или 0.
func CopyTimePointer(value *time.Time) *time.Time {
	if value == nil || value.IsZero() {
		return nil
	}

	c := *value

	return &c
}

func isNullable(required []bool) bool {
	return len(required) == 0 || !required[0]
}
