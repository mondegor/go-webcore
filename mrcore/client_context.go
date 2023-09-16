package mrcore

import (
    "context"
    "io"
    "net/http"
)

type (
    ClientData interface {
        Request() *http.Request
        RequestPath() RequestPath

        Context() context.Context
        WithContext(ctx context.Context) ClientData

        Writer() http.ResponseWriter

        Parse(structRequest any) error
        Validate(structRequest any) error
        ParseAndValidate(structRequest any) error

        SendResponse(status int, structResponse any) error
        SendResponseNoContent() error
        SendFile(contentType string, file io.Reader) error
    }

    RequestPath interface {
        Get(name string) string
        GetInt(name string) int64
    }
)
