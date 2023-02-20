package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j/db"
)

func ConvertMapStr(originalMap map[string]interface{}) map[string]string {
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

func PrettyPrintArray(arr []any) {
	b, err := json.MarshalIndent(arr, "", "  ")
	if err != nil {
		fmt.Println("marshal error:", err)
	}
	fmt.Print(string(b))
}

func PrettyPrintMapOfArrays(s map[string][]any) {
	b, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		fmt.Println("marshal error:", err)
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

func ConvertString(s string) any {
	// Try to convert to an int
	i, err := strconv.Atoi(s)
	if err == nil {
		return i
	}

	// Try to convert to a float
	f, err := strconv.ParseFloat(s, 64)
	if err == nil {
		return f
	}

	// Return the original string if it cannot be converted
	return s
}

func AnyToString(val interface{}) string {
	switch v := val.(type) {
	case int:
		return strconv.Itoa(v)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 64)
	}
	return ""
}

func RemoveElement(slice []string, i int) []string {
	if len(slice) == 0 || i < 0 || i >= len(slice) {
		// If the slice is empty or the index is out of range, return the original slice
		return slice
	}
	// Use append to create a new slice with the element at index i removed
	return append(slice[:i], slice[i+1:]...)
}

// not used right now. comparison happens in database. see query_builder_teimscaledb
// maybe change this or alternate depending on use case for performacne reasons
func CompareValues(val any, compareVal any, compareOp string) (bool, error) {
	fmt.Println("XXX CompareOperator: ", compareOp)
	fmt.Println("XXX VALUE: ", val)
	fmt.Println("XXX COMPAREVALUE: ", compareVal)
	switch v := val.(type) {
	case int:
		compareVal, ok := compareVal.(int)
		if !ok {
			return false, errors.New("error - compare value required to be an int")
		}
		switch compareOp {
		case "=":
			return v == compareVal, nil
		case ">":
			return v > compareVal, nil
		case "<":
			return v < compareVal, nil
		case ">=":
			return v >= compareVal, nil
		case "<=":
			return v <= compareVal, nil
		case "!=":
			return v != compareVal, nil
		default:
			return false, errors.New("error - compare operator not supported")
		}
	case float64:
		compareVal, ok := compareVal.(float64)
		if !ok {
			return false, errors.New("error - compare value required to be a float64")
		}
		switch compareOp {
		case "=":
			return v == compareVal, nil
		case ">":
			return v > compareVal, nil
		case "<":
			return v < compareVal, nil
		case ">=":
			return v >= compareVal, nil
		case "<=":
			return v <= compareVal, nil
		case "!=":
			return v != compareVal, nil
		default:
			return false, errors.New("error - compare operator not supported")
		}
	case float32:
		compareVal, ok := compareVal.(float32)
		if !ok {
			return false, errors.New("error - compare value required to be a float32")
		}
		switch compareOp {
		case "=":
			return v == compareVal, nil
		case ">":
			return v > compareVal, nil
		case "<":
			return v < compareVal, nil
		case ">=":
			return v >= compareVal, nil
		case "<=":
			return v <= compareVal, nil
		case "!=":
			return v != compareVal, nil
		default:
			return false, errors.New("error - compare operator not supported")
		}
	case string:
		compareVal, ok := compareVal.(string)
		// get rid of "" or ''. TODO: get the raw string without the quotation marks from the parse tree
		if strings.HasPrefix(compareVal, "\"") {
			compareVal = strings.Trim(compareVal, "\"")
		} else if strings.HasPrefix(compareVal, "'") {
			compareVal = strings.Trim(compareVal, "'")
		}
		if !ok {
			return false, errors.New("error - compare value required to be a string")
		}
		switch compareOp {
		case "=":
			return v == compareVal, nil
		case ">":
			return v > compareVal, nil
		case "<":
			return v < compareVal, nil
		case ">=":
			return v >= compareVal, nil
		case "<=":
			return v <= compareVal, nil
		case "!=":
			return v != compareVal, nil
		default:
			return false, errors.New("error - compare operator not supported")
		}
	default:
		return false, errors.New("error - unsupported value type")
	}
}

// used to disable annoying unused errors. delete the occourences in the end
func UNUSED(x ...interface{}) {}

func RemoveIdxFromSlice(slice []any, i int) []any {
	return append(slice[:i], slice[i+1:]...)
}
