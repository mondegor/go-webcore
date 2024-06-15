package mrflow

import (
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrstatus"
)

// ItemStatusFlow - возвращает стандартную карту возможных переходов ItemStatus.
func ItemStatusFlow() *StatusFlow {
	return NewStatusFlow(
		[]StatusFlowItem{
			{
				From: mrenum.ItemStatusDraft,
				To: []mrstatus.Getter{
					mrenum.ItemStatusEnabled,
					mrenum.ItemStatusDisabled,
				},
			},
			{
				From: mrenum.ItemStatusEnabled,
				To: []mrstatus.Getter{
					mrenum.ItemStatusDisabled,
				},
			},
			{
				From: mrenum.ItemStatusDisabled,
				To: []mrstatus.Getter{
					mrenum.ItemStatusEnabled,
				},
			},
		},
	)
}
