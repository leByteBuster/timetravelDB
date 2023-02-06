package main

import (
	"errors"
	"strings"
)

// parse the query string and return the start and end time, whether the query is shallow or not, and the cypher query string
func ParseQuery(query string) (string, string, bool, string, error) {

	// split the query by the space character
	split := strings.Split(query, " ")
	// the query has to have at least 8 elements: "FROM", $iso8601string1, "TO", $iso8601string2, "SHALLOW" or "MATCH..."
	if len(split) < 8 {
		return "", "", false, "", errors.New("invalid query")
	}
	// the first element has to be "FROM"
	if split[0] != "FROM" {
		return "", "", false, "", errors.New("invalid query")
	}
	// the third element has to be "TO"
	if split[2] != "TO" {
		return "", "", false, "", errors.New("invalid query")
	}
	// the fifth element has to be "SHALLOW" or "MATCH..."
	if split[4] != "SHALLOW" && split[4] != "MATCH" {
		return "", "", false, "", errors.New("invalid query")
	}
	// the sixth element has to be "MATCH" if the fifth element is "SHALLOW"
	if split[4] == "SHALLOW" && split[5] != "MATCH" {
		return "", "", false, "", errors.New("invalid query")
	}

	isShallow := false
	cypherQuery := ""
	if split[4] == "SHALLOW" {
		isShallow = true
		// put the rest of the query back together from the 7th element
		cypherQuery = strings.Join(split[5:], " ")
	} else {
		cypherQuery = strings.Join(split[4:], " ")
	}

	if !IsValidISO8601(split[1]) || !IsValidISO8601(split[3]) {
		return "", "", false, "", errors.New("invalid query")
	}

	// if IsValidCypher(cypherQuery) || IsValidISO8601(split[1]) || IsValidISO8601(split[3]) {
	// 	return "", "", false, "", errors.New("invalid query")
	// }

	return split[1], split[3], isShallow, cypherQuery, nil

}

// this function checks if the passed query string is a valid  Neo4j CYPHER string. For checking
// it uses the libcypher-parser of cleishm
func IsValidCypher(cypherQuery string) {
	panic("unimplemented")
}

// give me a function that checks if a string is a valid ISO8601 string
