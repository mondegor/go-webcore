package initing_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mondegor/go-core/errors"
	"github.com/mondegor/go-core/mraccess"
	"github.com/mondegor/go-core/mrlog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mondegor/go-webcore/mrcore/initing"
	"github.com/mondegor/go-webcore/mrserver"
)

type (
	// testRightsRegistry - реестр прав, в котором зарегистрированы все права, кроме notRegistered.
	testRightsRegistry struct {
		notRegistered string
	}
)

// IsRegistered - сообщает, что право зарегистрировано, если оно не совпадает с notRegistered.
func (r testRightsRegistry) IsRegistered(name string) bool {
	return name != r.notRegistered
}

func TestWithCheckAccessMiddleware_NilUserProvider(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name          string
		privilege     string
		permission    string
		notRegistered string
		accessToken   string
		wantErr       error
		wantCalled    bool
	}

	tests := []testCase{
		{
			name:       "public everyone",
			privilege:  mraccess.PrivilegePublic,
			permission: mraccess.PermissionEveryone,
			wantCalled: true,
		},
		{
			name:       "public guest-only",
			privilege:  mraccess.PrivilegePublic,
			permission: mraccess.PermissionGuestOnly,
			wantCalled: true,
		},
		{
			name:        "public guest-only with access token",
			privilege:   mraccess.PrivilegePublic,
			permission:  mraccess.PermissionGuestOnly,
			accessToken: "any-token",
			wantErr:     errors.ErrHttpAccessForbidden,
		},
		{
			name:       "public any-user",
			privilege:  mraccess.PrivilegePublic,
			permission: mraccess.PermissionAnyUser,
			wantErr:    errors.ErrHttpClientUnauthorized,
		},
		{
			name:       "privilege permission",
			privilege:  "admin",
			permission: "edit",
			wantErr:    errors.ErrHttpClientUnauthorized,
		},
		{
			name:       "privilege guest-only",
			privilege:  "admin",
			permission: mraccess.PermissionGuestOnly,
			wantErr:    errors.ErrHttpAccessForbidden,
		},
		{
			name:          "unregistered permission",
			privilege:     "admin",
			permission:    "edit",
			notRegistered: "edit",
			wantErr:       errors.ErrHttpAccessForbidden,
		},
	}

	internalHeaders := []string{
		mrserver.HeaderKeyInternalUserIDSlashGroup,
		mrserver.HeaderKeyInternalSessionID,
		mrserver.HeaderKeyInternalLangCode,
		mrserver.HeaderKeyInternalTimeZone,
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			called := false
			prepare := initing.WithCheckAccessMiddleware(
				mrlog.NopLogger(),
				mraccess.ActionGroup{Name: "test", Privilege: tt.privilege, BasePath: "/v1/"},
				nil,
				testRightsRegistry{notRegistered: tt.notRegistered},
			)

			handler := prepare(
				mrserver.HttpHandler{
					Method:     http.MethodGet,
					URL:        "/items",
					Permission: tt.permission,
					Func: func(_ http.ResponseWriter, r *http.Request) error {
						called = true

						// внутренние заголовки, подставленные клиентом, не должны дойти до обработчика
						for _, key := range internalHeaders {
							assert.Empty(t, r.Header.Values(key), "header %s must be removed", key)
						}

						return nil
					},
				},
			)

			assert.Equal(t, "/v1/items", handler.URL)

			r := httptest.NewRequest(http.MethodGet, "/v1/items", http.NoBody)

			for _, key := range internalHeaders {
				r.Header.Set(key, "spoofed-by-client")
			}

			if tt.accessToken != "" {
				r.Header.Set("Authorization", "Bearer "+tt.accessToken)
			}

			err := handler.Func(httptest.NewRecorder(), r)

			assert.Equal(t, tt.wantCalled, called)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
