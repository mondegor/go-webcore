package mrserver

import (
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type (
	requestPath struct {
		params httprouter.Params
	}
)

func newRequestPath(r *http.Request) *requestPath {
	params, ok := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)

	if !ok {
		params = nil
	}

	return &requestPath{
		params: params,
	}
}

func (r *requestPath) Get(name string) string {
	if r.params == nil {
		return ""
	}

	return r.params.ByName(name)
}

func (r *requestPath) GetInt64(name string) int64 {
	value, err := strconv.ParseInt(r.Get(name), 10, 64)

	if err != nil {
		return 0
	}

	return value
}
