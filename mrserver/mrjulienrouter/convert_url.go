package mrjulienrouter

import (
	"regexp"
	"strings"

	"github.com/mondegor/go-webcore/mrserver"
)

var regexpURLVars = regexp.MustCompile(`{([a-zA-Z][a-zA-Z0-9_]*)}`)

// ConvertURL - comment func.
func ConvertURL(url string) string {
	url = strings.Replace(url, mrserver.VarRestOfURL, "*"+varRestOfURL, 1)

	for _, m := range regexpURLVars.FindAllStringSubmatch(url, -1) {
		url = strings.Replace(url, m[0], ":"+m[1], 1)
	}

	return url
}
