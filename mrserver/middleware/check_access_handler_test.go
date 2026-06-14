package middleware_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mondegor/go-sysmess/mraccess"
	"github.com/mondegor/go-sysmess/mrlog"
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
			user := mraccess.NewUser([16]byte{1, 2, 3, 4}, "users", tc.sessionID, "", getter)
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
