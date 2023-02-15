package api

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j/db"
)

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

func prettyPrintMapOfArrays(s map[string][]any) {
	b, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Println(string(b))
}

// use this for large results from neo4j when the data
// have to be accessed multiple times
// Q: is this valid for every record of the result ? or do I have to run this function for every record anew
func InitIndex(rec *db.Record) map[string]int {
	keyMap := make(map[string]int, len(rec.Keys))
	for i, key := range rec.Keys {
		keyMap[key] = i
	}
	return keyMap
}

// TODO: test if faster
func GetIndexed(rec *db.Record, key string, keyMap map[string]int) (interface{}, bool) {
	i, ok := keyMap[key]
	if ok {
		return rec.Values[i], true
	}
	return nil, false
}

func UNUSED(x ...interface{}) {}
