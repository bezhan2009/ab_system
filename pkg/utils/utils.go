package utils

func GetString(m map[string]interface{}, key, def string) string {
	val, ok := m[key]
	if !ok {
		return def
	}
	str, ok := val.(string)
	if !ok {
		return def
	}
	return str
}
