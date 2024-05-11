package mrreq

import (
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/mondegor/go-webcore/mrcore"
)

const (
	maxLenUUID = 64
)

func ParseUUID(r *http.Request, key string, required bool) (uuid.UUID, error) {
	value := strings.TrimSpace(r.URL.Query().Get(key))

	if value == "" {
		if required {
			return uuid.Nil, mrcore.FactoryErrHTTPRequestParamEmpty.New(key)
		}

		return uuid.Nil, nil
	}

	if len(value) > maxLenUUID {
		return uuid.Nil, mrcore.FactoryErrHTTPRequestParamLenMax.New(key, maxLenInt64)
	}

	item, err := uuid.Parse(value)
	if err != nil {
		return uuid.Nil, mrcore.FactoryErrHTTPRequestParseParam.New(key, "UUID", value)
	}

	return item, nil
}
