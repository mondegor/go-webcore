package parser

import (
	"net/http"

	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-sysmess/mrtype/parse"
)

type (
	// EnumList - парсер перечислений в виде строки.
	EnumList[T ~uint8] struct {
		logger       mrlog.Logger
		defaultItems []T
		parser       func(items []string) ([]T, error)
	}
)

// NewEnumList - создаёт объект  EnumList.
func NewEnumList[T ~uint8](
	logger mrlog.Logger,
	parser func(items []string) ([]T, error),
) *EnumList[T] {
	return &EnumList[T]{
		logger: logger,
		parser: parser,
	}
}

// NewEnumListWithDefault - создаёт объект EnumList со статусами по умолчанию.
func NewEnumListWithDefault[T ~uint8](
	logger mrlog.Logger,
	items []T,
	parser func(items []string) ([]T, error),
) *EnumList[T] {
	return &EnumList[T]{
		logger:       logger,
		defaultItems: items,
		parser:       parser,
	}
}

// FilterEnumList - возвращает массив mrenum. EnumList поступивший из внешнего запроса.
// Если ключ key не найден или возникнет ошибка, то возвращается значение по умолчанию.
func (p *EnumList[T]) FilterEnumList(r *http.Request, key string) []T {
	items, err := p.parseList(r, key)
	if err != nil {
		p.logger.Warn(r.Context(), "FilterEnumList", "key", key, "error", err)

		return p.defaultItems
	}

	if len(items) == 0 {
		return p.defaultItems
	}

	return items
}

func (p *EnumList[T]) parseList(r *http.Request, key string) ([]T, error) {
	enumList, err := parse.EnumList(r.URL.Query().Get(key))
	if err != nil {
		return nil, err
	}

	items, err := p.parser(enumList)
	if err != nil {
		return nil, err
	}

	return items, nil
}
