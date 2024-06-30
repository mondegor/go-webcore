package mrview

import (
	"regexp"
)

var (
	regexpAnyNotSpaceSymbol = regexp.MustCompile(`^\S+$`)
	regexpArticle           = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_.+-]*[a-zA-Z0-9]$`)
	regexpDoubleSize        = regexp.MustCompile(`^[0-9]+x[0-9]+$`)
	regexpRewriteName       = regexp.MustCompile(`^[a-z][a-z0-9-]*[a-z0-9]$`)
	regexpVariable          = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9]+$`)
)

// ValidateAnyNotSpaceSymbol - comment func.
func ValidateAnyNotSpaceSymbol(value string) bool {
	return regexpAnyNotSpaceSymbol.MatchString(value)
}

// ValidateArticle - comment func.
func ValidateArticle(value string) bool {
	return regexpArticle.MatchString(value)
}

// ValidateDoubleSize - comment func.
func ValidateDoubleSize(value string) bool {
	return regexpDoubleSize.MatchString(value)
}

// ValidateRewriteName - comment func.
func ValidateRewriteName(value string) bool {
	return regexpRewriteName.MatchString(value)
}

// ValidateVariable - comment func.
func ValidateVariable(value string) bool {
	return regexpVariable.MatchString(value)
}
