package mail_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mondegor/go-webcore/mrclient"
	"github.com/mondegor/go-webcore/mrclient/mail"
)

// Make sure the mail.SMTPClient conforms with the mrclient.MailSender interface.
func TestSMTPClientImplementsMailSender(t *testing.T) {
	assert.Implements(t, (*mrclient.MailSender)(nil), &mail.SMTPClient{})
}
