package mrjson_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/mrjson"
)

// Make sure the JsonEncoder conforms with the mrserver.ResponseEncoder interface.
func TestJsonEncoderImplementsResponseEncoder(t *testing.T) {
	assert.Implements(t, (*mrserver.ResponseEncoder)(nil), &mrjson.JsonEncoder{})
}
