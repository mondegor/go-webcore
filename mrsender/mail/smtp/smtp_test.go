package smtp_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mondegor/go-webcore/mrsender"
	"github.com/mondegor/go-webcore/mrsender/mail/smtp"
)

// Make sure the smtp.MailClient conforms with the mrsender.MailProvider interface.
func TestMailClientImplementsMailProvider(t *testing.T) {
	assert.Implements(t, (*mrsender.MailProvider)(nil), &smtp.MailClient{})
}
