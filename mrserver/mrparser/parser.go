package mrparser

type (
	// Parser - агрегатор базовых парсеров.
	Parser struct {
		*Bool
		*DateTime
		*Int64
		*Uint64
		*ListSorter
		*ListPager
		*String
		*UUID
		*Validator
		*File
		*Image
	}
)

// NewParser - создаёт объект Parser.
func NewParser(
	p1 *Bool,
	p2 *DateTime,
	p3 *Int64,
	p4 *Uint64,
	p5 *ListSorter,
	p6 *ListPager,
	p7 *String,
	p8 *UUID,
	p9 *Validator,
	p10 *File,
	p11 *Image,
) *Parser {
	return &Parser{
		Bool:       p1,
		DateTime:   p2,
		Int64:      p3,
		Uint64:     p4,
		ListSorter: p5,
		ListPager:  p6,
		String:     p7,
		UUID:       p8,
		Validator:  p9,
		File:       p10,
		Image:      p11,
	}
}
