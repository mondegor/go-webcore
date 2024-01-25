package mrtype

func StringPointer(value string) *string {
	if value == "" {
		return nil
	}

	return &value
}

func StringPointerCopy(value *string) *string {
	if value == nil || *value == "" {
		return nil
	}

	c := *value

	return &c
}
