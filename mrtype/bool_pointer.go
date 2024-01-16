package mrtype

var (
	boolFalse = newFalse()
	boolTrue  = newTrue()
)

func BoolPointer(value bool) *bool {
	if value {
		return boolTrue
	}

	return boolFalse
}

func newFalse() *bool {
	value := false
	return &value
}

func newTrue() *bool {
	value := true
	return &value
}
