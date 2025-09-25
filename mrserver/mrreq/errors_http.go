package mrreq

import (
	"github.com/mondegor/go-sysmess/mrerr"
)

var (
	// ErrHttpRequestCorrelationID - header 'X-Correlation-Id' contains incorrect value.
	// Это вспомогательная ошибка, для неё отключено формирование стека вызовов и отправление события о её создании.
	ErrHttpRequestCorrelationID = mrerr.NewKindInternal("header 'X-Correlation-Id' contains incorrect value: '{Value}'",
		mrerr.WithDisabledCaller(), mrerr.WithDisabledOnCreated())

	// ErrHttpRequestParseClientIP - clientIP contains incorrect value.
	// Это вспомогательная ошибка, для неё отключено формирование стека вызовов и отправление события о её создании.
	ErrHttpRequestParseClientIP = mrerr.NewKindInternal("clientIP contains incorrect value: '{Value}'",
		mrerr.WithDisabledCaller(), mrerr.WithDisabledOnCreated())
)
