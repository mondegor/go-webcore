package mrresp_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mondegor/go-webcore/mrserver/mrresp"
)

// TestNewError400Response_TimeFormat - фиксирует формат поля Time: RFC3339 с точностью
// до миллисекунд и суффиксом Z. Клиенты разбирают поле строгим форматом, поэтому потеря
// дробной части (напр. возврат к time.RFC3339) должна ломать тест, а не клиента.
func TestNewError400Response_TimeFormat(t *testing.T) {
	t.Parallel()

	const layout = "2006-01-02T15:04:05.000Z07:00"

	before := time.Now().UTC().Truncate(time.Millisecond)

	r := httptest.NewRequest(http.MethodPost, "/orders", http.NoBody)
	resp := mrresp.NewError400Response(r)

	after := time.Now().UTC()

	require.NotEmpty(t, resp.Time)
	assert.Equal(t, http.StatusBadRequest, resp.Status)
	assert.Equal(t, "POST /orders", resp.Instance)

	// строгий разбор тем же layout: значение, записанное другим форматом, здесь упадёт
	parsed, err := time.Parse(layout, resp.Time)
	require.NoError(t, err, "поле Time должно быть записано форматом %q, получено %q", layout, resp.Time)

	assert.True(t, strings.HasSuffix(resp.Time, "Z"), "время формируется в UTC, ожидается суффикс Z: %q", resp.Time)
	assert.Len(t, resp.Time, len("2006-01-02T15:04:05.000Z"))

	// время должно попадать в интервал вызова
	assert.False(t, parsed.Before(before), "Time=%s раньше начала вызова %s", parsed, before)
	assert.False(t, parsed.After(after), "Time=%s позже конца вызова %s", parsed, after)
}
