package api

func getNullableString(input map[string]interface{}, key string) string {
	s := ""
	r := input[key]
	if r != nil {
		s = r.(string)
	}
	return s
}

func getNullOrString(input map[string]interface{}, key string) (string, bool) {
	r := input[key]
	if r != nil {
		return r.(string), true
	}
	return "", false
}

func getNullableInt(input map[string]interface{}, key string) int {
	var s int
	r := input[key]
	if r != nil {
		s = r.(int)
	}
	return s
}

func getFalseOrBoolean(input map[string]interface{}, key string) bool {
	givenValue := input[key]
	if givenValue != nil {
		return givenValue.(bool)
	}
	return false
}
