package mrjson_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mondegor/go-webcore/mrserver/mrjson"
	"github.com/mondegor/go-webcore/mrserver/request"
)

// Make sure the JsonDecoder conforms with the request.ParserDecode interface.
func TestJsonDecoderImplementsRequestDecoder(t *testing.T) {
	t.Parallel()

	assert.Implements(t, (*request.ParserDecode)(nil), &mrjson.JsonDecoder{})
}
