package mrinit

import (
	"fmt"
	"io"

	"github.com/mondegor/go-sysmess/mrerr/mr"
	"github.com/mondegor/go-sysmess/mrlog"
)

type (
	// ResourceList - comment struct.
	ResourceList struct {
		logger mrlog.LiteLogger
		list   []io.Closer
	}
)

// NewResourceList - создаёт объект ResourceList.
func NewResourceList(logger mrlog.LiteLogger) *ResourceList {
	return &ResourceList{
		logger: logger,
	}
}

// Add - последовательно вызываются все функции из указанного списка.
func (rl *ResourceList) Add(resource io.Closer) {
	rl.list = append(rl.list, resource)
}

// Close - последовательно вызываются все функции из указанного списка.
func (rl *ResourceList) Close() {
	for _, resource := range rl.list {
		rl.close(resource)
	}
}

func (rl *ResourceList) close(resource io.Closer) {
	if err := resource.Close(); err != nil {
		rl.logger.Error(
			"ResourceList.Close()",
			"error", mr.ErrInternalFailedToClose.Wrap(err),
			"resource", fmt.Sprintf("%#v", resource),
		)

		return
	}

	rl.logger.Info(fmt.Sprintf("Resource %T was closed", resource))
}
