package mrjson_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/mrjson"
)

// Make sure the JsonDecoder conforms with the mrserver.RequestDecoder interface.
func TestJsonDecoderImplementsRequestDecoder(t *testing.T) {
	assert.Implements(t, (*mrserver.RequestDecoder)(nil), &mrjson.JsonDecoder{})
}
