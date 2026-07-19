package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mondegor/go-core/errors"
	"github.com/mondegor/go-core/mrlog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/middleware"
)

// TestCheckAccessTokenHandler_InternalHeadersAreDropped - проверяет, что внутренние
// заголовки, подставленные клиентом, не доходят до обработчика: иначе клиент мог бы
// выдать себя за авторизованного пользователя, задав их вручную.
func TestCheckAccessTokenHandler_InternalHeadersAreDropped(t *testing.T) {
	t.Parallel()

	internalHeaders := []string{
		mrserver.HeaderKeyUserIDSlashGroup,
		mrserver.HeaderKeyTimeZone,
		mrserver.HeaderKeySessionID,
	}

	var got http.Header

	next := func(_ http.ResponseWriter, r *http.Request) error {
		got = r.Header.Clone()

		return nil
	}

	handler := middleware.CheckAccessTokenHandler(mrlog.NopLogger(), "test")(next)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", http.NoBody)

	for _, key := range internalHeaders {
		r.Header.Set(key, "spoofed-by-client")
	}

	// заголовок клиента, не являющийся внутренним, должен остаться нетронутым
	r.Header.Set(mrserver.HeaderKeyAcceptLanguage, "en")

	require.NoError(t, handler(w, r))

	for _, key := range internalHeaders {
		// проверяется отсутствие записи в карте, а не пустое значение:
		// заголовок должен быть удалён, а не обнулён
		_, ok := got[http.CanonicalHeaderKey(key)]
		assert.False(t, ok, "header %s must be removed", key)
	}

	assert.Equal(t, "en", got.Get(mrserver.HeaderKeyAcceptLanguage))
}

// TestCheckAccessTokenHandler_AuthorizedIsForbidden - проверяет, что запрос
// с access token до обработчика не доходит.
func TestCheckAccessTokenHandler_AuthorizedIsForbidden(t *testing.T) {
	t.Parallel()

	nextCalled := false

	next := func(_ http.ResponseWriter, _ *http.Request) error {
		nextCalled = true

		return nil
	}

	handler := middleware.CheckAccessTokenHandler(mrlog.NopLogger(), "test")(next)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	r.Header.Set("Authorization", "Bearer any-token")

	require.ErrorIs(t, handler(w, r), errors.ErrHttpAccessForbidden)
	assert.False(t, nextCalled)
}
