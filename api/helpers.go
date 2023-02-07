package api

import "time"

func convertMapStr(originalMap map[string]interface{}) map[string]string {
	convertedMap := map[string]string{}
	for key, value := range originalMap {
		convertedMap[key] = value.(string)
	}

	return convertedMap
}

func IsValidISO8601(s string) bool {
	_, err := time.Parse("2006-01-02T15:04:05.9999999999Z", s)
	if err != nil {
		return false
	}
	return true
}
