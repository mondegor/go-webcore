package parser

import (
	"net/http"

	"github.com/mondegor/go-core/mrlog"
	"github.com/mondegor/go-core/mrtype"
	"github.com/mondegor/go-core/mrtype/parse"
)

type (
	// Float64 - парсер числа с плавающей запятой.
	Float64 struct {
		logger mrlog.Logger
	}
)

// NewFloat64 - создаёт объект Float64.
func NewFloat64(logger mrlog.Logger) *Float64 {
	return &Float64{
		logger: logger,
	}
}

// FilterFloat64 - возвращает число с плавающей запятой поступившее из внешнего запроса.
// Если ключ key не найден или возникнет ошибка, то возвращается нулевое значение.
func (p *Float64) FilterFloat64(r *http.Request, key string) float64 {
	value, err := parse.Float64(r.URL.Query().Get(key), false)
	if err != nil {
		p.logger.Warn(r.Context(), "FilterFloat64", "key", key, "error", err)

		return 0
	}

	return value
}

// FilterRangeFloat64 - возвращает интервал состоящий из двух чисел с плавающей запятой поступивший из внешнего запроса.
// Если ключ key с постфиксами "-min", "-max" не найдены или возникнет ошибка, то возвращается нулевой интервал.
func (p *Float64) FilterRangeFloat64(r *http.Request, key string) mrtype.RangeFloat64 {
	value, err := parse.RangeFloat64(r.URL.Query().Get(key+"-min"), r.URL.Query().Get(key+"-max"))
	if err != nil {
		p.logger.Warn(r.Context(), "FilterRangeFloat64", "key", key, "error", err)

		return mrtype.RangeFloat64{}
	}

	return value
}
