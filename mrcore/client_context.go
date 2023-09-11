package mrcore

import (
    "context"
    "net/http"
)

type (
    ClientData interface {
        Request() *http.Request
        Context() context.Context
        WithContext(ctx context.Context) ClientData
        RequestPath() RequestPath
        Writer() http.ResponseWriter

        Parse(structRequest any) error
        Validate(structRequest any) error
        ParseAndValidate(structRequest any) error

        SendResponse(status int, structResponse any) error
        SendResponseNoContent() error
    }

    RequestPath interface {
        Get(name string) string
        GetInt(name string) int64
    }
)
