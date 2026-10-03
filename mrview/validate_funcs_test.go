package mrview_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mondegor/go-webcore/mrview"
)

// TestValidatePassword - пароль допускает любые печатные символы ASCII, кроме пробела.
func TestValidatePassword(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{name: "letters and digits", value: "Abc123xyz", want: true},
		{name: "all special chars", value: "!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~", want: true},
		{name: "backtick", value: "pass`word", want: true},
		{name: "empty", value: "", want: false},
		{name: "space", value: "pass word", want: false},
		{name: "tab", value: "pass\tword", want: false},
		{name: "cyrillic", value: "пароль123", want: false},
		{name: "del char", value: "pass\x7fword", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, mrview.ValidatePassword(tt.value))
		})
	}
}
