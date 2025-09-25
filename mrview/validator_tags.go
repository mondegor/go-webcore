package mrview

import (
	"regexp"
)

var (
	regexpAnyNotSpaceSymbol = regexp.MustCompile(`^\S+$`)
	regexpVariable          = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9]*$`)
	regexpName              = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9/_.+-]*[a-zA-Z0-9]$`)
	regexpRewriteName       = regexp.MustCompile(`^[a-z][a-z0-9-]*[a-z0-9]$`)
	regexpPassword          = regexp.MustCompile(`^[a-zA-Z0-9!"#$%&'()*+,\-./:;<=>?@\[\\\]^_{|}~]+$`)
	regexpDoubleSize        = regexp.MustCompile(`^[0-9]+x[0-9]+$`)
	regexpTripleSize        = regexp.MustCompile(`^[0-9]+x[0-9]+x[0-9]+$`)
)

// ValidateAnd - comment func.
func ValidateAnd(values ...func(value string) bool) func(value string) bool {
	return func(value string) bool {
		for _, fn := range values {
			if !fn(value) {
				return false
			}
		}

		return true
	}
}

// ValidateOr - comment func.
func ValidateOr(values ...func(value string) bool) func(value string) bool {
	return func(value string) bool {
		for _, fn := range values {
			if fn(value) {
				return true
			}
		}

		return false
	}
}

// ValidateAnyNotSpaceSymbol - comment func.
func ValidateAnyNotSpaceSymbol(value string) bool {
	return regexpAnyNotSpaceSymbol.MatchString(value)
}

// ValidateVariable - comment func.
func ValidateVariable(value string) bool {
	return regexpVariable.MatchString(value)
}

// ValidateName - comment func.
func ValidateName(value string) bool {
	return regexpName.MatchString(value)
}

// ValidateRewriteName - comment func.
func ValidateRewriteName(value string) bool {
	return regexpRewriteName.MatchString(value)
}

// ValidatePassword - comment func.
func ValidatePassword(value string) bool {
	return regexpPassword.MatchString(value)
}

// ValidateDoubleSize - comment func.
func ValidateDoubleSize(value string) bool {
	return regexpDoubleSize.MatchString(value)
}

// ValidateTripleSize - comment func.
func ValidateTripleSize(value string) bool {
	return regexpTripleSize.MatchString(value)
}
