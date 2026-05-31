package telegram_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mondegor/go-webcore/mrclient"
	"github.com/mondegor/go-webcore/mrclient/telegram"
)

// Make sure the telegram.BotClient conforms with the mrclient.MessengerSender interface.
func TestMessageClientImplementsMessageProvider(t *testing.T) {
	t.Parallel()

	assert.Implements(t, (*mrclient.MessengerSender)(nil), &telegram.BotClient{})
}
