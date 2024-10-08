package mrreqinit

import (
	"github.com/mondegor/go-webcore/mrcore/mrinit"
	"github.com/mondegor/go-webcore/mrserver/mrreq"
)

// EnrichedHttpErrors - возвращает список ошибок наделённых дополнительными
// параметрами используемые в http обработчиках.
func EnrichedHttpErrors() []mrinit.EnrichedError {
	return []mrinit.EnrichedError{
		mrinit.WrapProtoExtraDisabled(mrreq.ErrHttpRequestCorrelationID),
		mrinit.WrapProtoExtraDisabled(mrreq.ErrHttpRequestUserIP),
		mrinit.WrapProtoExtraDisabled(mrreq.ErrHttpRequestParseUserIP),
	}
}
