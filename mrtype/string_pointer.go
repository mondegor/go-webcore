package mrtype

func StringPointer(value string) *string {
	if value == "" {
		return nil
	}

	return &value
}
