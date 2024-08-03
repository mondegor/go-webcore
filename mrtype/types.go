package mrtype

import "github.com/mondegor/go-webcore/mrenum"

type (
	// KeyInt32 - целочисленный ID записи.
	KeyInt32 int32

	// KeyInt64 - целочисленный ID записи.
	KeyInt64 int64

	// KeyString - строковый ID записи.
	KeyString string

	// RangeInt64 - целочисленный интервал [Min, Max].
	RangeInt64 struct {
		Min int64
		Max int64
	}

	// RangeFloat64 - вещественный интервал [Min, Max].
	RangeFloat64 struct {
		Min float64
		Max float64
	}

	// PageParams - параметры для выборки части списка элементов.
	PageParams struct {
		Index uint64 // pageIndex
		Size  uint64 // pageSize
	}

	// SortParams - параметры для сортировки списка элементов по указанному полю.
	SortParams struct {
		FieldName string               // sortField
		Direction mrenum.SortDirection // sortDirection
	}
)

// Transform - преобразовывает в RangeFloat64 с умножением полей на указанный коэффициент
// (для приведения к необходимой ед. измерения).
func (r RangeInt64) Transform(coefficient float64) RangeFloat64 {
	return RangeFloat64{
		Min: float64(r.Min) * coefficient,
		Max: float64(r.Max) * coefficient,
	}
}
