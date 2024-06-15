package mrreqinit

import (
	"github.com/mondegor/go-webcore/mrcore/mrinit"
	"github.com/mondegor/go-webcore/mrserver/mrreq"
)

// ManagedHttpErrors - comment func.
func ManagedHttpErrors() []mrinit.ManagedError {
	return []mrinit.ManagedError{
		mrinit.WrapProtoExtraDisabled(mrreq.ErrHttpRequestCorrelationID),
		mrinit.WrapProtoExtraDisabled(mrreq.ErrHttpRequestUserIP),
		mrinit.WrapProtoExtraDisabled(mrreq.ErrHttpRequestParseUserIP),
	}
}
