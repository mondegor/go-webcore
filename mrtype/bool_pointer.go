package mrtype

func BoolPointer(value bool) *bool {
	return &value
}

func BoolPointerCopy(value *bool) *bool {
	if value == nil {
		return nil
	}

	c := *value

	return &c
}
