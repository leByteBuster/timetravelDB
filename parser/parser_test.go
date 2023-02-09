package parser

import (
	"testing"
)

func TestParseQueryValid(t *testing.T) {
	valid_queries := []string{
		"FROM 2022-12-22T15:33:13.4Z TO 2022-12-29T20:24:36.311107Z SHALLOW MATCH (n:Node) RETURN n",
		"FROM 2022-12-22T15:33:13Z TO 2022-12-29T20:24:36.311107Z MATCH (n:Node) RETURN n",
		"FROM 2022-12-22T15:33:13Z TO 2022-12-29T20:24:36.311107Z MATCH (n:Node) WHERE xyz RETURN n"}
	for _, query := range valid_queries {
		_, err := ParseQuery(query)
		if err != nil {
			t.Fatalf("Query should be valid: %v\n", query)
		}
	}

}

func TestParseQueryInvalid(t *testing.T) {
	invalid_queries := []string{
		"FROM 2022-12-22T:33:13Z TO 2022-12-29T20:24:36.311107Z SHALLOW MATCH (n:Node) RETURN n",
		"FROM 2022-12-22T15:33:13Z TO 2022-12-29T20:24:36.311107GMT MATCH (n:Node) RETURN n",
		"FROM 2022-12-22T15:33:13Z TO 2022-12-29T20:24:36.311107Z  (n:Node) WHERE xyz RETURN n"}
	for _, query := range invalid_queries {
		_, err := ParseQuery(query)
		if err == nil {
			t.Fatalf("Query should be invalid: %v\n", query)
		}
	}
}
