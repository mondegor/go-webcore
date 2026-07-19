package middleware_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mondegor/go-core/mraccess"
	"github.com/mondegor/go-core/mrlog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/middleware"
)

type (
	// testRightsSource - источник прав ролей для построения RightsGetter в тестах.
	testRightsSource map[string][]string

	// spyLogger - логгер, фиксирующий вызовы уровня Error.
	spyLogger struct {
		mrlog.Logger
		errCount int
		lastMsg  string
	}
)

// RoleRights - возвращает права указанной роли и признак её наличия.
func (s testRightsSource) RoleRights(role string) (rights []string, ok bool) {
	rights, ok = s[role]

	return rights, ok
}

// Error - фиксирует факт логирования ошибки.
func (l *spyLogger) Error(_ context.Context, msg string, _ ...any) {
	l.errCount++
	l.lastMsg = msg
}

func newTestRightsGetter(t *testing.T) mraccess.RightsGetter {
	t.Helper()

	getter, err := mraccess.NewRolesGroupSet(
		[]mraccess.RoleGroup{{Name: "users", Roles: []string{"all"}}},
		testRightsSource{"all": {mraccess.PermissionEveryone}},
	)
	require.NoError(t, err)

	return getter
}

func TestCheckAccessHandler_UserLangAndTimeZone(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name         string
		userLang     string
		userTimeZone string
		clientLang   string
		clientTZ     string
		wantLang     string
		wantTZ       string
	}

	tests := []testCase{
		{
			name:         "user values are set",
			userLang:     "ru",
			userTimeZone: "Europe/Moscow",
			wantLang:     "ru",
			wantTZ:       "Europe/Moscow",
		},
		{
			name:         "client values are overwritten by user values",
			userLang:     "ru",
			userTimeZone: "Europe/Moscow",
			clientLang:   "en",
			clientTZ:     "Asia/Tokyo",
			wantLang:     "ru",
			wantTZ:       "Europe/Moscow",
		},
		{
			// Accept-Language - клиентский заголовок, при пустом значении в профиле
			// он не удаляется и доходит до обработчика как есть, а внутренний
			// заголовок с часовым поясом срезается, чтобы клиент не мог его подделать
			name:       "client lang survives, but client time zone is removed",
			clientLang: "en",
			clientTZ:   "Asia/Tokyo",
			wantLang:   "en",
			wantTZ:     "",
		},
		{
			name:     "headers are absent when neither profile nor client set them",
			wantLang: "",
			wantTZ:   "",
		},
	}

	getter := newTestRightsGetter(t)
	action := mraccess.Action{
		Privilege:  mraccess.PrivilegePublic,
		Permission: mraccess.PermissionEveryone,
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			logger := &spyLogger{Logger: mrlog.NopLogger()}
			user := mraccess.NewUser([16]byte{1, 2, 3, 4}, "users", "8f14e45f", tc.userLang, tc.userTimeZone, getter)
			provider := mraccess.NewOneUserProvider(user)

			var got http.Header

			next := func(_ http.ResponseWriter, r *http.Request) error {
				got = r.Header.Clone()

				return nil
			}

			handler := middleware.CheckAccessHandler(logger, action, provider)(next)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
			r.Header.Set("Authorization", "Bearer any-token")

			if tc.clientLang != "" {
				r.Header.Set(mrserver.HeaderKeyAcceptLanguage, tc.clientLang)
			}

			if tc.clientTZ != "" {
				r.Header.Set(mrserver.HeaderKeyTimeZone, tc.clientTZ)
			}

			require.NoError(t, handler(w, r))

			// пустое ожидаемое значение означает отсутствие самой записи в карте
			// заголовков, а не запись с пустым значением
			for key, want := range map[string]string{
				mrserver.HeaderKeyAcceptLanguage: tc.wantLang,
				mrserver.HeaderKeyTimeZone:       tc.wantTZ,
			} {
				_, ok := got[http.CanonicalHeaderKey(key)]

				if want == "" {
					assert.False(t, ok, "header %s must be removed", key)

					continue
				}

				assert.True(t, ok, "header %s must be set", key)
				assert.Equal(t, want, got.Get(key))
			}
		})
	}
}

func TestCheckAccessHandler_SessionID(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name         string
		sessionID    string
		wantErrCount int
	}

	tests := []testCase{
		{
			name:         "session id is set",
			sessionID:    "8f14e45f",
			wantErrCount: 0,
		},
		{
			name:         "session id is empty",
			sessionID:    "",
			wantErrCount: 1,
		},
	}

	getter := newTestRightsGetter(t)
	action := mraccess.Action{
		Privilege:  mraccess.PrivilegePublic,
		Permission: mraccess.PermissionEveryone,
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			logger := &spyLogger{Logger: mrlog.NopLogger()}
			user := mraccess.NewUser([16]byte{1, 2, 3, 4}, "users", tc.sessionID, "", "", getter)
			provider := mraccess.NewOneUserProvider(user)

			var (
				nextCalled   bool
				gotSessionID string
				hasHeader    bool
			)

			next := func(_ http.ResponseWriter, r *http.Request) error {
				nextCalled = true
				gotSessionID = r.Header.Get(mrserver.HeaderKeySessionID)
				_, hasHeader = r.Header[http.CanonicalHeaderKey(mrserver.HeaderKeySessionID)]

				return nil
			}

			handler := middleware.CheckAccessHandler(logger, action, provider)(next)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
			r.Header.Set("Authorization", "Bearer any-token")

			require.NoError(t, handler(w, r))

			assert.True(t, nextCalled)
			assert.True(t, hasHeader)
			assert.Equal(t, tc.sessionID, gotSessionID)
			assert.Equal(t, tc.wantErrCount, logger.errCount)

			if tc.wantErrCount > 0 {
				assert.Equal(t, "session id is empty", logger.lastMsg)
			}
		})
	}
}
