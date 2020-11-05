package elasticsearch

func mapStringArray(source []interface{}) []string {
	result := []string{}
	for _, item := range source {
		result = append(result, item.(string))
	}
	return result
}

func mapStringMap(source map[string]interface{}) map[string]string {
	var result = make(map[string]string)

	for key, value := range source {
		result[key] = value.(string)
	}
	return result
}
