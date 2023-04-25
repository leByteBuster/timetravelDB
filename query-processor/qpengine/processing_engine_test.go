package qpengine

import (
	"log"
	"strings"
	"testing"

	"github.com/LexaTRex/timetravelDB/query-processor/parser"
	"github.com/LexaTRex/timetravelDB/utils"
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
		utils.Debugf("MatchClause: %v", res.MatchClause)
		utils.Debugf("WhereClause: %v", res.WhereClause)
		utils.Debugf("ReturnClause: %v", res.ReturnClause)
		tmpQuery := buildTmpWhereClause(res.From, res.To, res.WhereClause, res.QueryVariables.MatchQueryVariables)

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
		"FROM 2123-12-13T12:34:39Z TO 2123-12-13T14:34:39.2222Z MATCH (n) WHERE n.ts_ping > 22.33 RETURN n ",
		"FROM 2123-12-13T12:34:39Z TO 2123-12-13T14:34:39.2222Z MATCH (n) WHERE n.ping > 22.33 RETURN n ",
		"FROM 2123-12-13T12:34:39Z TO 2123-12-13T14:34:39.2222Z MATCH (n) WHERE n.ts_ping > 22.33 AND n.ts_name = 'hans' RETURN n ",
		"FROM 2123-12-13T12:34:39Z TO 2123-12-13T14:34:39.2222Z MATCH (n) WHERE n.ping > 22.33 AND n.name = 'hans' RETURN n ",
		"FROM 2123-12-13T12:34:39Z TO 2123-12-13T14:34:39.2222Z MATCH (n) WHERE n.ts_ping > 22.33 OR n.ts_ping IS NOT NULL RETURN n ",
		"FROM 2123-12-13T12:34:39Z TO 2123-12-13T14:34:39.2222Z MATCH (n) WHERE n.ping > 22.33 OR n.ping IS NOT NULL RETURN n ",
		"FROM 2123-12-13T12:34:39Z TO 2123-12-13T14:34:39.2222Z MATCH (n) WHERE n.ts_ping IS NOT NULL RETURN n ",
		"FROM 2123-12-13T12:34:39Z TO 2123-12-13T14:34:39.2222Z MATCH (n) WHERE n.ping IS NOT NULL RETURN n ",
		// note: OR is not supported yet
		"FROM 2123-12-13T12:34:39Z TO 2123-12-13T14:34:39.2222Z MATCH (n) WHERE n.ts_ping > 22.33 AND n.ts_ping < 23 OR n.ts_ping IS NULL RETURN n ",
	}
	expected_results := []string{
		"WHERE n.ts_ping IS NOT NULL",
		"WHERE n.ping > 22.33",
		"WHERE n.ts_ping IS NOT NULL AND n.ts_name IS NOT NULL",
		"WHERE n.ping > 22.33 AND n.name = 'hans'",
		"WHERE n.ts_ping IS NOT NULL OR n.ts_ping IS NOT NULL",
		"WHERE n.ping > 22.33 OR n.ping IS NOT NULL",
		"WHERE n.ts_ping IS NOT NULL",
		"WHERE n.ping IS NOT NULL",
		"WHERE n.ts_ping IS NOT NULL AND n.ts_ping IS NOT NULL OR n.ts_ping IS NULL",
		"WHERE n.ping > 22.33 AND n.ping < 23 OR n.ping IS NULL",
	}
	invalid_queries := []string{
		"FROM 2123-12-13T12:34:39Z TO 2123-12-13T14:34:39.2222Z MATCH (n) WHERE n.ping > 22.33 > 11 RETURN n ",
		"FROM 2123-12-13T12:34:39Z TO 2123-12-13T14:34:39.2222Z MATCH (n) WHERE n.ping > n.ding > 3 RETURN n ",
		"FROM 2123-12-13T12:34:39Z TO 2123-12-13T14:34:39.2222Z MATCH (n) WHERE n.ping > n.ping NOT NULL RETURN n ",                      // maybe before error
		"FROM 2123-12-13T12:34:39Z TO 2123-12-13T14:34:39.2222Z MATCH (n) WHERE n.ping > 20 AND (n.ping > 10 OR n.ping > 30)  RETURN n ", // this is still valid but should be invalid
	}
	utils.UNUSED(invalid_queries)
	for i, query := range valid_queries {
		res, err := parser.ParseQuery(query)
		if err != nil {
			t.Fatalf("Parsing error: %v\n This should not be happening with valid queries\n", query)
		}
		whereClause := res.WhereClause
		manipulated, err := buildCondWhereClause(res.LookupsWhereRelevant, whereClause)
		if err != nil {
			t.Fatalf("Manipulation error for: %v\n Manipulated WHERE clause is not as expected\n", query)
		}
		if strings.Trim(manipulated, " ") != strings.Trim(expected_results[i], " ") {
			t.Fatalf("\nManipulation error for:\n     %v\nManipulated WHERE clause:\n     %v\nis not as expected:\n     %v\n", query, manipulated, expected_results[i])
		}
		log.Println(whereClause)
	}

	// TODO: update test with invalid queries
	// 			 there is limited possibility to check for validity as long as the parser is passed so we are not protected against invalid queries
	// for i, query := range invalid_queries {
	// 	res, err1 := parser.ParseQuery(query)
	// 	if err1 != nil {
	// 		whereClause := res.WhereClause
	// 		manipulated, err := buildCondWhereClause(res.LookupsWhereRelevant, whereClause)
	// 		if err != nil {
	// 			if strings.Trim(manipulated, " ") == strings.Trim(expected_results[i], " ") {
	// 				t.Fatalf("\nThis should be invalid: %v\n", query)
	// 			}
	// 		}
	// 	}
	// }
}
