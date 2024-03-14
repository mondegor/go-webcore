package mrparser

type (
	Parser struct {
		*Bool
		*DateTime
		*Int64
		*KeyInt32
		*ListSorter
		*ListPager
		*String
		*UUID
		*Validator
		*File
		*Image
	}
)

func NewParser(
	p1 *Bool,
	p2 *DateTime,
	p3 *Int64,
	p4 *KeyInt32,
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
		KeyInt32:   p4,
		ListSorter: p5,
		ListPager:  p6,
		String:     p7,
		UUID:       p8,
		Validator:  p9,
		File:       p10,
		Image:      p11,
	}
}
