package dataadapter

import (
	"encoding/json"
	"os"
)

func LoadJsonData(path string) ([]map[string]interface{}, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Decode the JSON data
	var data []map[string]interface{}
	err = json.NewDecoder(file).Decode(&data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func ConvertMaps(originalMaps []interface{}) []map[string]interface{} {
	convertedMaps := make([]map[string]interface{}, 0)
	for _, originalMap := range originalMaps {
		convertedMap := map[string]interface{}{}
		for key, value := range originalMap.(map[string]interface{}) {
			convertedMap[key] = value.(interface{})
		}
		convertedMaps = append(convertedMaps, convertedMap)
	}
	return convertedMaps
}
