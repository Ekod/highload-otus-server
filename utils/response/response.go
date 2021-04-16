package response

//Data формирует ответ клиенту
func Data(key string, value interface{}) map[string]interface{} {
	responseMap := make(map[string]interface{})
	responseMap[key] = value

	return responseMap
}
