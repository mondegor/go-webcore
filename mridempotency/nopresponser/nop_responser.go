package nopresponser

import (
	"github.com/mondegor/go-webcore/mridempotency"
)

type (
	// Responser - comment struct.
	Responser struct{}
)

// Make sure the Image conforms with the mridempotency.Responser interface.
var _ mridempotency.Responser = (*Responser)(nil)

// New - создаёт объект Responser.
func New() *Responser {
	return &Responser{}
}

// StatusCode - comment method.
func (r Responser) StatusCode() int {
	return 0
}

// Body - comment method.
func (r Responser) Body() []byte {
	return nil
}
