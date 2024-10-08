package nopprovider_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mondegor/go-webcore/mridempotency"
	"github.com/mondegor/go-webcore/mridempotency/nopprovider"
)

// Make sure the Provider conforms with the mridempotency.Provider interface.
func TestProviderImplementsProvider(t *testing.T) {
	assert.Implements(t, (*mridempotency.Provider)(nil), &nopprovider.Provider{})
}
