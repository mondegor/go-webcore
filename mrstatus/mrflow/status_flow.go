package mrflow

import "github.com/mondegor/go-webcore/mrstatus"

type (
	// StatusFlow - comment struct.
	StatusFlow struct {
		fromToMap     map[mrstatus.Getter][]mrstatus.Getter
		toFromMap     map[mrstatus.Getter][]mrstatus.Getter
		registeredMap map[mrstatus.Getter]bool
	}

	// StatusFlowItem - comment struct.
	StatusFlowItem struct {
		From mrstatus.Getter
		To   []mrstatus.Getter
	}
)

// NewStatusFlow - создаёт объект StatusFlow.
func NewStatusFlow(list []StatusFlowItem) *StatusFlow {
	fromToMap := make(map[mrstatus.Getter][]mrstatus.Getter, len(list))
	toFromMap := make(map[mrstatus.Getter][]mrstatus.Getter, len(list))
	registeredMap := make(map[mrstatus.Getter]bool, len(list))

	for _, item := range list {
		fromToMap[item.From] = item.To
		registeredMap[item.From] = true

		for _, to := range item.To {
			toFromMap[to] = append(toFromMap[to], item.From)
			registeredMap[to] = true
		}
	}

	return &StatusFlow{
		fromToMap: fromToMap,
		toFromMap: toFromMap,
	}
}

// Exists - проверяет, что данный статус зарегистрирован в карте статусов.
func (f *StatusFlow) Exists(status mrstatus.Getter) bool {
	_, ok := f.registeredMap[status]

	return ok
}

// Check - проверяет возможность переключения из указанного статуса в указанный статус.
func (f *StatusFlow) Check(from, to mrstatus.Getter) bool {
	toStatuses, ok := f.fromToMap[from]
	if !ok {
		return false
	}

	for i := range toStatuses {
		if toStatuses[i] == to {
			return true
		}
	}

	return false
}

// PossibleToStatuses - возвращает список статусов в которые можно переключить указанный статус.
func (f *StatusFlow) PossibleToStatuses(from mrstatus.Getter) []mrstatus.Getter {
	if toStatuses, ok := f.fromToMap[from]; ok {
		return toStatuses
	}

	return nil
}

// PossibleFromStatuses - возвращается список статусов из которых можно переключиться в указанный статус.
func (f *StatusFlow) PossibleFromStatuses(to mrstatus.Getter) []mrstatus.Getter {
	if fromStatuses, ok := f.toFromMap[to]; !ok {
		return fromStatuses
	}

	return nil
}
