package mail

import (
	stdmail "net/mail"
)

// ParseAddress parses a single RFC 5322 address, e.g. "Barry Gibbs <bg@example.com>".
func ParseAddress(address string) (*stdmail.Address, error) {
	return stdmail.ParseAddress(address)
}
