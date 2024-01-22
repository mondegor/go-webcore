package mrview

import (
	"regexp"
)

var (
	regexpAnyNotSpaceSymbol = regexp.MustCompile(`^\S+$`)
	regexpArticle           = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_.+-]*[a-zA-Z0-9]$`)
	regexpVariable          = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9]+$`)
)

func ValidateAnyNotSpaceSymbol(value string) bool {
	return regexpAnyNotSpaceSymbol.MatchString(value)
}

func ValidateArticle(value string) bool {
	return regexpArticle.MatchString(value)
}

func ValidateVariable(value string) bool {
	return regexpVariable.MatchString(value)
}
