package mrresp_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/mrresp"
)

// Make sure the Sender conforms with the mrserver.ResponseSender interface.
func TestSenderImplementsResponseSender(t *testing.T) {
	assert.Implements(t, (*mrserver.ResponseSender)(nil), &mrresp.Sender{})
}
