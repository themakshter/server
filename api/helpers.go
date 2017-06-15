package api

func getNullableString(input map[string]interface{}, key string) string {
	s := ""
	r := input[key]
	if r != nil {
		s = r.(string)
	}
	return s
}

func getNullableInt(input map[string]interface{}, key string) int {
	var s int
	r := input[key]
	if r != nil {
		s = r.(int)
	}
	return s
}
