package mrresp_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/mrresp"
)

// Make sure the ErrorSender conforms with the mrserver.ErrorResponseSender interface.
func TestErrorSenderImplementsErrorResponseSender(t *testing.T) {
	assert.Implements(t, (*mrserver.ErrorResponseSender)(nil), &mrresp.ErrorSender{})
}
