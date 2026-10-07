package validate_test

import (
	"testing"

	"github.com/mondegor/go-core/mrlog"
	"github.com/stretchr/testify/assert"

	"github.com/mondegor/go-webcore/mrserver/request/parser"
	"github.com/mondegor/go-webcore/mrserver/request/validate"
)

// Проверка на этапе компиляции: при встраивании парсеров методы не конфликтуют.
var (
	_ validate.RequestParser        = (*validate.Parser)(nil)
	_ validate.RequestListParser    = (*validate.ListParser)(nil)
	_ validate.RequestContextParser = (*validate.ContextParser)(nil)
)

// TestNewParser - базовый парсер хранит переданные парсеры.
func TestNewParser(t *testing.T) {
	t.Parallel()

	boolParser := parser.NewBool(mrlog.NopLogger())
	float64Parser := parser.NewFloat64(mrlog.NopLogger())

	got := validate.NewParser(boolParser, nil, nil, float64Parser, nil, nil, nil, nil)

	assert.Same(t, boolParser, got.Bool)
	assert.Same(t, float64Parser, got.Float64)
	assert.Nil(t, got.String)
}

// TestNewListParser - парсер параметров списка хранит переданные парсеры.
func TestNewListParser(t *testing.T) {
	t.Parallel()

	listPagerParser := parser.NewListPager(mrlog.NopLogger(), parser.ListPagerOptions{})
	listSorterParser := parser.NewListSorter(mrlog.NopLogger(), parser.ListSorterOptions{})

	got := validate.NewListParser(listPagerParser, listSorterParser)

	assert.Same(t, listPagerParser, got.ListPager)
	assert.Same(t, listSorterParser, got.ListSorter)

	got = validate.NewListParser(nil, listSorterParser)

	assert.Nil(t, got.ListPager)
	assert.Same(t, listSorterParser, got.ListSorter)
}

// TestNewContextParser - парсер контекста запроса хранит переданные парсеры.
func TestNewContextParser(t *testing.T) {
	t.Parallel()

	clientParser := parser.NewClient(mrlog.NopLogger(), parser.ClientOptions{})
	userParser := parser.NewUser(mrlog.NopLogger())

	got := validate.NewContextParser(clientParser, userParser, nil, nil)

	assert.Same(t, clientParser, got.Client)
	assert.Same(t, userParser, got.User)
	assert.Nil(t, got.Locale)
	assert.Nil(t, got.TimeZone)
}
