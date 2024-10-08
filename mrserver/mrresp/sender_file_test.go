package mrresp_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/mrresp"
)

// Make sure the FileSender conforms with the mrserver.FileResponseSender interface.
func TestFileSenderImplementsFileResponseSender(t *testing.T) {
	assert.Implements(t, (*mrserver.FileResponseSender)(nil), &mrresp.FileSender{})
}
