package api

import (
	"log"
	"testing"

	"github.com/LexaTRex/timetravelDB/parser"
)

func TestParseQueryValid(t *testing.T) {
	valid_queries := []string{
		"FROM 2123-12-13T12:34:39Z TO 2123-12-13T14:34:39.2222Z MATCH (n) WHERE n.ping > 22.33 RETURN n.ping, n ",
		"FROM 2023-02-03T12:34:39Z TO 2023-02-03T14:34:39.2222Z SHALLOW MATCH (a)-[x]->(b) RETURN a.ping, b ",
		"FROM 2023-02-03T12:34:39Z TO 2023-02-03 SHALLOW MATCH (a)-[x]->(b) WHERE a.ping > 22 RETURN a.ping, b ",
		"FROM 2023-02-03T12:34:39Z TO 2023-02-03 SHALLOW MATCH (a)-[x]->(b) WHERE a.ping > 22 RETURN a ",
		"FROM 2023-02-03T12:34:39Z TO 2023-02-03 SHALLOW MATCH (a) WHERE a.ping > 22 RETURN a ",
		"FROM 2123-12-13T12:34:39Z TO 2123-12-13T14:34:39.2222Z MATCH (a)-[x]->(b) WHERE a.ping > 22 RETURN a.ping, b ",
		"FROM 2023-02-03T12:34:39Z TO 2023-02-03 SHALLOW MATCH (a)-[x]->(b) WHERE a.ping > 22 RETURN a.ping, b ",
		"FROM 2023-02-03T12:34:39Z TO 2023-02-03 SHALLOW MATCH (a)-[x]->(b) WHERE a.ping > 22 RETURN a ",
		"FROM 2023-02-03T12:34:39Z TO 2023-02-03 SHALLOW MATCH (a) WHERE a.ping > 22 RETURN a ",
		"FROM 2123-12-13T12:34:39Z TO 2123-12-13T14:34:39.2222Z MATCH (a)-[x]->(b) WHERE a.ping > 22 RETURN a.ping, b ",
		"FROM 2022-12-22T15:33:13.4Z TO 2022-12-29T20:24:36.311107Z SHALLOW MATCH (n:Node) RETURN n",
		"FROM 2022-12-22T15:33:13Z TO 2022-12-29T20:24:36.311107Z MATCH (n:Node) RETURN n",
		"FROM 2022-12-22T15:33:13Z TO 2022-12-29T20:24:36.311107Z MATCH (n:Node) WHERE xyz RETURN n"}
	for _, query := range valid_queries {
		res, err := parser.ParseQuery(query)
		log.Printf("MatchClause: %v", res.MatchClause)
		log.Printf("WhereClause: %v", res.WhereClause)
		log.Printf("ReturnClause: %v", res.ReturnClause)
		tmpQuery := addTempToWhereQuery(res.From, res.To, res.WhereClause, res.GraphElements.MatchGraphElements)

		log.Println(tmpQuery)
		if err != nil {
			t.Fatalf("Query should be valid: %v\n", query)
		}
		// TODO: comparing against expected values
		// if tmpQuery != expected {
		// 	t.Fatalf("Query should be valid: %v\n", query)
		// }

	}

}

// func TestParseQueryShallowIsNullInWherePropInReturn(t *testing.T) {
// 	query := "FROM 2123-12-13T12:34:39Z TO 2123-12-13T14:34:39.2222Z SHALLOW MATCH (a)-[x]->(b) WHERE a.ping IS NOT NULL" + " RETURN a.ping, b "
// 	res, err := parser.ParseQuery(query)
// }

func TestManipulateWhereClauseNeo4j(t *testing.T) {
	valid_queries := []string{
		"FROM 2123-12-13T12:34:39Z TO 2123-12-13T14:34:39.2222Z MATCH (n) WHERE n.ping > 22.33 RETURN n ",
		"FROM 2123-12-13T12:34:39Z TO 2123-12-13T14:34:39.2222Z MATCH (n) WHERE n.ping > 22.33 AND n.name = 'hans' RETURN n ",
		"FROM 2123-12-13T12:34:39Z TO 2123-12-13T14:34:39.2222Z MATCH (n) WHERE n.ping > 22.33 OR n.ping IS NOT NULL RETURN n ",
		"FROM 2123-12-13T12:34:39Z TO 2123-12-13T14:34:39.2222Z MATCH (n) WHERE n.ping IS NOT NULL RETURN n ",
		"FROM 2123-12-13T12:34:39Z TO 2123-12-13T14:34:39.2222Z MATCH (n) WHERE n.ping > 22.33 AND n.ping < 23 OR n.ping IS NULL RETURN n ",
	}
	expected_results := []string{
		"WHERE n.ping IS NOT NULL",
		"WHERE n.ping IS NOT NULL AND n.name IS NOT NULL",
		"WHERE n.ping IS NOT NULL OR n.ping IS NOT NULL",
		"WHERE n.ping IS NOT NULL",
		"WHERE n.ping IS NOT NULL AND n.ping IS NOT NULL OR n.ping IS NULL",
	}
	invalid_queries := []string{
		"FROM 2123-12-13T12:34:39Z TO 2123-12-13T14:34:39.2222Z MATCH (n) WHERE n.ping > 22.33 > 11 RETURN n ",
		"FROM 2123-12-13T12:34:39Z TO 2123-12-13T14:34:39.2222Z MATCH (n) WHERE n.ping > n.ding > 3 RETURN n ",
		"FROM 2123-12-13T12:34:39Z TO 2123-12-13T14:34:39.2222Z MATCH (n) WHERE n.ping > n.ping NOT NULL RETURN n ",                      // maybe before error
		"FROM 2123-12-13T12:34:39Z TO 2123-12-13T14:34:39.2222Z MATCH (n) WHERE n.ping > 20 AND (n.ping > 10 OR n.ping > 30)  RETURN n ", // this is still valid but should be invalid
	}
	for i, query := range valid_queries {
		res, err := parser.ParseQuery(query)
		if err != nil {
			t.Fatalf("Parsing error: %v\n This should not be happening with valid queries\n", query)
		}
		whereClause := res.WhereClause
		manipulated, err := manipulateWhereClause(res, whereClause)
		if err != nil {
			t.Fatalf("Manipulation error for: %v\n Manipulated WHERE clause is not as expected\n", query)
		}
		if manipulated != expected_results[i] {
			t.Fatalf("Manipulation error for:\n     %v\nManipulated WHERE clause:\n     %v\nis not as expected:\n     %v\n", query, manipulated, expected_results[i])
		}
		log.Println(whereClause)
	}

	for i, query := range invalid_queries {
		res, err1 := parser.ParseQuery(query)
		if err1 != nil {
			whereClause := res.WhereClause
			manipulated, err := manipulateWhereClause(res, whereClause)
			if err != nil {
				if manipulated == expected_results[i] {
					t.Fatalf("\nThis should be invalid: %v\n", query)
				}
			}
		}
	}
}

// TODO: move this test to Parser where it belongs
func TestGetRelevantLookupInfoWhere(t *testing.T) {
	validQueries := []string{
		"FROM 2023-02-03T12:34:39Z TO 2023-02-03 SHALLOW MATCH (a)-[x]->(b) WHERE 22 > a.ping RETURN a",
		"FROM 2023-02-03T12:34:39Z TO 2023-02-03 SHALLOW MATCH (a)-[x]->(b) WHERE a.ping > 22 RETURN a",
		"FROM 2021-12-22T15:33:13.0000005Z TO 2024-01-12T15:33:13.0000006Z SHALLOW MATCH (a)-[x]->(b) WHERE a.name = 'UGWJn' RETURN  *",
	}
	expected := []parser.LookupInfo{{ElementVariable: "a", Property: "ping", CompareOperator: ">", CompareValue: 22, LookupLeft: false},
		{ElementVariable: "a", Property: "ping", CompareOperator: ">", CompareValue: 22, LookupLeft: true},
		{ElementVariable: "a", Property: "name", CompareOperator: "=", CompareValue: "'UGWJn'", LookupLeft: true},
	}

	for i, query := range validQueries {
		res, err := parser.ParseQuery(query)
		if err != nil {
			t.Fatalf("Query should be valid: %v\n", query)
		}
		lookups, err := parser.GetRelevantLookupInfoWhere(res)
		lookup := lookups[0]
		if err != nil {
			t.Fatalf("Error retrieving lookup info: %v\n", query)
		}
		if lookup.ElementVariable != expected[i].ElementVariable || lookup.Property != expected[i].Property || lookup.CompareOperator != expected[i].CompareOperator || lookup.CompareValue != expected[i].CompareValue || lookup.LookupLeft != expected[i].LookupLeft {
			t.Fatalf("\n Expected:\n%+v\nGot:\n%+v\n", expected[i], lookup)
		}
	}
}
