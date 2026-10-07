package validate

import (
	"github.com/mondegor/go-webcore/mrserver/request"
	"github.com/mondegor/go-webcore/mrserver/request/parser"
)

type (
	// RequestListParser - агрегирующий интерфейс парсеров параметров списка HTTP-запроса:
	// постраничная навигация и сортировка.
	RequestListParser interface {
		request.ParserListPager
		request.ParserListSorter
	}

	// ListParser - реализация RequestListParser на основе парсеров go-webcore.
	ListParser struct {
		*parser.ListPager
		*parser.ListSorter
	}
)

// NewListParser - создаёт объект ListParser.
func NewListParser(
	listPagerParser *parser.ListPager,
	listSorterParser *parser.ListSorter,
) *ListParser {
	return &ListParser{
		ListPager:  listPagerParser,
		ListSorter: listSorterParser,
	}
}
