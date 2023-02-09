package api

import (
	"fmt"

	"github.com/LexaTRex/timetravelDB/parser"
)

// process a TTQL Query
func ProcessQuery(query string) error {
	res, err := parser.ParseQuery(query)
	if err != nil {
		return err
	}

	fmt.Printf("res: %v", res)
	return nil
}

// this function checks if the passed query string is a valid  Neo4j CYPHER string. For checking
// it uses the libcypher-parser of cleishm
func IsValidCypher(cypherQuery string) {
	panic("unimplemented")
}

// give me a function that checks if a string is a valid ISO8601 string
