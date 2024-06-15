package mrjulienrouter_test

import (
	"testing"

	"github.com/go-playground/assert/v2"

	"github.com/mondegor/go-webcore/mrserver/mrjulienrouter"
)

func TestConvertURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		url  string
		want string
	}{
		{
			name: "test 1",
			url:  "/v1/prefix",
			want: "/v1/prefix",
		},
		{
			name: "test 2",
			url:  "{id}",
			want: ":id",
		},
		{
			name: "test 3",
			url:  "/v1/prefix/{id}",
			want: "/v1/prefix/:id",
		},
		{
			name: "test 4",
			url:  "{id}/postfix",
			want: ":id/postfix",
		},
		{
			name: "test 5",
			url:  "/v1/prefix/{id}/postfix",
			want: "/v1/prefix/:id/postfix",
		},
		{
			name: "Multi 1",
			url:  "{id1}/middle/{id2}",
			want: ":id1/middle/:id2",
		},
		{
			name: "Multi 2",
			url:  "/v1/prefix/{id1}/middle/{id2}",
			want: "/v1/prefix/:id1/middle/:id2",
		},
		{
			name: "Multi 3",
			url:  "/v1/prefix/{id1}/middle/{id2}/postfix",
			want: "/v1/prefix/:id1/middle/:id2/postfix",
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := mrjulienrouter.ConvertURL(tt.url)
			assert.Equal(t, tt.want, got)
		})
	}
}
