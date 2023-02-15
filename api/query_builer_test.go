package api

import (
	"log"
	"testing"

	"github.com/LexaTRex/timetravelDB/parser"
)

func TestParseQueryValid(t *testing.T) {
	valid_queries := []string{
		"FROM 2123-12-13T12:34:39Z TO 2123-12-13T14:34:39.2222Z MATCH (n) WHERE n.ping > 22.33" + "RETURN n.ping, n ",
		"FROM 2023-02-03T12:34:39Z TO 2023-02-03T14:34:39.2222Z SHALLOW MATCH (a)-[x]->(b) " + "RETURN a.ping, b ",
		"FROM 2023-02-03T12:34:39Z TO 2023-02-03 SHALLOW MATCH (a)-[x]->(b) WHERE a.ping > 22" + " RETURN a.ping, b ",
		"FROM 2023-02-03T12:34:39Z TO 2023-02-03 SHALLOW MATCH (a)-[x]->(b) WHERE a.ping > 22" + " RETURN a ",
		"FROM 2023-02-03T12:34:39Z TO 2023-02-03 SHALLOW MATCH (a) WHERE a.ping > 22" + " RETURN a ",
		"FROM 2123-12-13T12:34:39Z TO 2123-12-13T14:34:39.2222Z MATCH (a)-[x]->(b) WHERE a.ping > 22" + " RETURN a.ping, b ",
		"FROM 2023-02-03T12:34:39Z TO 2023-02-03 SHALLOW MATCH (a)-[x]->(b) WHERE a.ping > 22" + " RETURN a.ping, b ",
		"FROM 2023-02-03T12:34:39Z TO 2023-02-03 SHALLOW MATCH (a)-[x]->(b) WHERE a.ping > 22" + " RETURN a ",
		"FROM 2023-02-03T12:34:39Z TO 2023-02-03 SHALLOW MATCH (a) WHERE a.ping > 22" + " RETURN a ",
		"FROM 2123-12-13T12:34:39Z TO 2123-12-13T14:34:39.2222Z MATCH (a)-[x]->(b) WHERE a.ping > 22" + " RETURN a.ping, b ",
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
