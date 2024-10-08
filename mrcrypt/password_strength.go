package mrcrypt

import (
	"math"
)

const (
	PassStrengthNotRated PassStrength = iota // PassStrengthNotRated - пароль без оценки
	PassStrengthWeak                         // PassStrengthWeak - слабый пароль
	PassStrengthMedium                       // PassStrengthMedium - средний пароль
	PassStrengthStrong                       // PassStrengthStrong - надёжный пароль
	PassStrengthBest                         // PassStrengthBest - самый надёжный пароль
)

const (
	passTypeSmallABC passTypeChars = iota
	passTypeBigABC
	passTypeNumeral
	passTypeSign
)

const (
	passTypeNumeralLen  = 10
	passTypeBigABCLen   = 26
	passTypeSmallABCLen = 26
	passTypeSignLen     = 20
)

type (
	// PassStrength - надёжность пароля.
	PassStrength uint8

	passTypeChars uint8
)

var passStrengthName = map[PassStrength]string{ //nolint:gochecknoglobals
	PassStrengthNotRated: "NOT_RATED",
	PassStrengthWeak:     "WEAK",
	PassStrengthMedium:   "MIDDLE",
	PassStrengthStrong:   "STRONG",
	PassStrengthBest:     "THE_BEST",
}

// PasswordStrength - comment func.
func (l *Lib) PasswordStrength(value string) PassStrength {
	length := len(value)

	if length == 0 {
		return PassStrengthNotRated
	}

	uniqChars := make(map[byte]bool, length)
	uniqTypeChars := make(map[passTypeChars]int)

	for i := 0; i < length; i++ {
		uniqChars[value[i]] = true

		switch {
		case value[i] >= 48 && value[i] <= 57:
			uniqTypeChars[passTypeNumeral] = passTypeNumeralLen
		case value[i] >= 65 && value[i] <= 90:
			uniqTypeChars[passTypeBigABC] = passTypeBigABCLen
		case value[i] >= 97 && value[i] <= 122:
			uniqTypeChars[passTypeSmallABC] = passTypeSmallABCLen
		default:
			uniqTypeChars[passTypeSign] = passTypeSignLen
		}
	}

	var totalLen int

	for i := range uniqTypeChars {
		totalLen += uniqTypeChars[i]
	}

	if len(uniqTypeChars) > 1 { // минимально два набора символов должно использоваться
		// вычисление информационной энтропии
		bits := int(float64(len(uniqChars)) * math.Log(float64(totalLen)) / math.Ln2)

		if bits >= 67 { // min(11 uniq chars and 4 sets[69] OR 12 uniq chars and 3 sets[69] OR 13 uniq chars and 2 sets[67])
			return PassStrengthBest
		}

		if bits >= 56 { // min(9 uniq chars and 4 sets[57] OR 10 uniq chars and 3 sets[58] OR 11 uniq chars and 2 sets[56])
			return PassStrengthStrong
		}

		if bits >= 44 { // min(7 uniq chars and 4 sets[44] OR 8 uniq chars and 3 sets[46] OR 9 uniq chars and 2 sets[46])
			return PassStrengthMedium
		}

		if bits >= 31 { // min(5 uniq chars and 4 sets[31] OR 6 uniq chars and 3 sets[34] OR 7 uniq chars and 2 sets[36])
			return PassStrengthWeak
		}
	}

	return PassStrengthNotRated
}

// String - comment method.
func (e PassStrength) String() string {
	return passStrengthName[e]
}
