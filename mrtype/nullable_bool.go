package mrtype

const (
	NullableBoolNull  NullableBool = -1
	NullableBoolTrue               = 1
	NullableBoolFalse              = 0
)

type (
	NullableBool int8
)

func (b NullableBool) Val() bool {
	return b > 0
}

func (b *NullableBool) Set(value bool) {
	if value {
		*b = NullableBoolTrue
	} else {
		*b = NullableBoolFalse
	}
}

func (b NullableBool) IsNull() bool {
	return b < 0
}

func (b *NullableBool) SetNull() {
	*b = NullableBoolNull
}
