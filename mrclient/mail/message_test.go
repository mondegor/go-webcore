package mail_test

import (
	"mime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mondegor/go-webcore/mrclient/mail"
)

func TestNewMessage_ContentType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		opts    []mail.MessageOption
		want    string
		wantErr bool
	}{
		{
			name: "option is absent",
			want: "text/plain; charset=UTF-8",
		},
		{
			name: "empty value",
			opts: []mail.MessageOption{mail.WithContentType("")},
			want: "text/plain; charset=UTF-8",
		},
		{
			name: "html without charset",
			opts: []mail.MessageOption{mail.WithContentType("text/html")},
			want: "text/html; charset=UTF-8",
		},
		{
			name: "html with utf-8 charset",
			opts: []mail.MessageOption{mail.WithContentType("text/html; charset=utf-8")},
			want: "text/html; charset=UTF-8",
		},
		{
			name: "upper case type and charset param",
			opts: []mail.MessageOption{mail.WithContentType("TEXT/HTML; Charset=UTF-8")},
			want: "text/html; charset=UTF-8",
		},
		{
			name: "rfc 2231 encoded param is kept",
			opts: []mail.MessageOption{mail.WithContentType("text/plain; name*=utf-8''%D1%82")},
			want: "text/plain; charset=UTF-8; name*=utf-8''%D1%82",
		},
		{
			name: "other params are kept",
			opts: []mail.MessageOption{mail.WithContentType("text/plain; format=flowed; charset=utf-8")},
			want: "text/plain; charset=UTF-8; format=flowed",
		},
		{
			name: "quoted param value is kept quoted",
			opts: []mail.MessageOption{mail.WithContentType(`text/plain; name="a b"`)},
			want: `text/plain; charset=UTF-8; name="a b"`,
		},
		{
			name: "multipart without charset",
			opts: []mail.MessageOption{mail.WithContentType("Multipart/Alternative; boundary=xyz")},
			want: "multipart/alternative; boundary=xyz",
		},
		{
			name: "non-text type keeps its params",
			opts: []mail.MessageOption{mail.WithContentType("application/json; charset=utf-8")},
			want: "application/json; charset=utf-8",
		},
		{
			name:    "invalid params",
			opts:    []mail.MessageOption{mail.WithContentType("text/html;;")},
			wantErr: true,
		},
		{
			name:    "text with another charset",
			opts:    []mail.MessageOption{mail.WithContentType("TEXT/HTML; charset=koi8-r")},
			wantErr: true,
		},
		{
			name:    "upper case charset param name with another charset",
			opts:    []mail.MessageOption{mail.WithContentType("text/html; Charset=KOI8-R")},
			wantErr: true,
		},
		{
			name:    "media type without subtype",
			opts:    []mail.MessageOption{mail.WithContentType("text")},
			wantErr: true,
		},
		{
			name:    "invalid media type",
			opts:    []mail.MessageOption{mail.WithContentType("text/")},
			wantErr: true,
		},
		{
			name:    "wildcard subtype",
			opts:    []mail.MessageOption{mail.WithContentType("text/*")},
			wantErr: true,
		},
		{
			name:    "wildcard media type",
			opts:    []mail.MessageOption{mail.WithContentType("*/*")},
			wantErr: true,
		},
		{
			name:    "header injection via CRLF",
			opts:    []mail.MessageOption{mail.WithContentType("text/plain\r\nBcc: evil@example.com")},
			wantErr: true,
		},
		{
			name:    "header injection via CRLF in quoted param",
			opts:    []mail.MessageOption{mail.WithContentType("text/plain; name=\"a\r\nBcc: evil@example.com\"")},
			wantErr: true,
		},
		{
			name:    "header injection via rfc 2231 encoded CRLF",
			opts:    []mail.MessageOption{mail.WithContentType("text/plain; name*=utf-8''%0D%0ABcc:evil@example.com")},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			msg, err := mail.NewMessage("from@example.com", "to@example.com", tt.opts...)

			if tt.wantErr {
				require.ErrorIs(t, err, mail.ErrInternalParsingContentTypeFailed)
				assert.Nil(t, msg)

				return
			}

			require.NoError(t, err)

			got := msg.Header().Get("Content-Type")
			assert.Equal(t, tt.want, got)

			_, _, err = mime.ParseMediaType(got)
			assert.NoError(t, err)
		})
	}
}
