package mrcore

//go:generate mockgen -source=id_generator.go -destination=./mock/id_generator.go

type (
	// IdentifierGenerator - comment interface.
	IdentifierGenerator interface {
		GenID() string
	}

	// IdentifierGeneratorFunc - comment func type.
	IdentifierGeneratorFunc func() string
)

// GenID - comment method.
func (f IdentifierGeneratorFunc) GenID() string {
	return f()
}
