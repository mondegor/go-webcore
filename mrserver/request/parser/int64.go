package parser

import (
	"net/http"

	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-sysmess/mrtype"
	"github.com/mondegor/go-sysmess/mrtype/parse"
)

type (
	// Int64 - парсер int числа.
	Int64 struct {
		logger mrlog.Logger
	}
)

// NewInt64 - создаёт объект Int64.
func NewInt64(
	logger mrlog.Logger,
) *Int64 {
	return &Int64{
		logger: logger,
	}
}

// FilterInt64 - возвращает int64 число поступившее из внешнего запроса.
// Если ключ key не найден или возникнет ошибка, то возвращается нулевое значение.
func (p *Int64) FilterInt64(r *http.Request, key string) int64 {
	value, err := parse.Int64(r.URL.Query().Get(key), false)
	if err != nil {
		p.logger.Warn(r.Context(), "FilterInt64", "key", key, "error", err)

		return 0
	}

	return value
}

// FilterRangeInt64 - возвращает интервал состоящий из двух int64 чисел поступивший из внешнего запроса.
// Если ключ key с постфиксами "-min", "-max" не найдены или возникнет ошибка, то возвращается нулевой интервал.
func (p *Int64) FilterRangeInt64(r *http.Request, key string) mrtype.RangeInt64 {
	value, err := parse.RangeInt64(r.URL.Query().Get(key+"-min"), r.URL.Query().Get(key+"-max"))
	if err != nil {
		p.logger.Warn(r.Context(), "FilterRangeInt64", "key", key, "error", err)

		return mrtype.RangeInt64{}
	}

	return value
}

// FilterInt64List - возвращает массив int64 чисел поступивший из внешнего запроса.
// Если ключ key не найден или возникнет ошибка, то возвращается nil значение.
func (p *Int64) FilterInt64List(r *http.Request, key string) []int64 {
	items, err := parse.Int64List(r.URL.Query().Get(key))
	if err != nil {
		p.logger.Warn(r.Context(), "FilterInt64List", "key", key, "error", err)

		return nil
	}

	return items
}
