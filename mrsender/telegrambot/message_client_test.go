package telegrambot_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mondegor/go-webcore/mrsender"
	"github.com/mondegor/go-webcore/mrsender/telegrambot"
)

// Make sure the telegrambot.MessageClient conforms with the mrsender.MessageProvider interface.
func TestMessageClientImplementsMessageProvider(t *testing.T) {
	assert.Implements(t, (*mrsender.MessageProvider)(nil), &telegrambot.MessageClient{})
}
