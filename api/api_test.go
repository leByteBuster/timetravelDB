package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"

	databaseapi "github.com/LexaTRex/timetravelDB/database-api"
	"github.com/LexaTRex/timetravelDB/parser"
	"github.com/LexaTRex/timetravelDB/utils"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

/**
  * This is a test integration test file.
  * It tests the API by sending requests to the API and checking the response.
  * For the tests to work, a Neo4j database is required to be running on the default port,
  * as well as a timetravelDB database running on the default port. The databases are required
	* to have the same credentials as the ones in the tests. Furthermore they are required to be initialized
	* with the according testing data stored in /test_data
  * See DOCKER_README.md for setting up the testing environment.
**/

var testConfNeo = databaseapi.Neo4jConfig{
	Host:     "localhost",
	Port:     "7687",
	Username: "neo4j",
	Password: "test",
}

var testConfTS = databaseapi.TimescaleConfig{
	Host:     "localhost",
	Port:     "5432",
	Username: "postgres",
	Password: "password",
	Database: "postgres",
}

func TestDeepQueries(t *testing.T) {

	// initialize Neo4j
	DriverNeo, err := neo4j.NewDriverWithContext("neo4j://"+testConfNeo.Host+":"+testConfNeo.Port, neo4j.BasicAuth(testConfNeo.Username, testConfNeo.Password, ""))
	if err != nil {
		log.Printf("Creating driver failed: %v", err)
		os.Exit(1)
	}
	defer DriverNeo.Close(context.Background())

	databaseapi.SessionNeo = DriverNeo.NewSession(context.Background(), neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer databaseapi.SessionNeo.Close(context.Background())

	// initialize TimescaleDB
	databaseapi.SessionTS, err = databaseapi.ConnectTimescale(testConfTS.Username, testConfTS.Password, testConfTS.Port, testConfTS.Database)
	if err != nil {
		log.Printf("Creating driver failed: %v", err)
		os.Exit(1)
	}
	defer databaseapi.SessionTS.Close()

	query1 := "FROM 2021-12-22T15:33:13.0000005Z TO 2024-01-12T15:33:13.0000006Z  MATCH (a)-[x]->(b) RETURN  a,x,b"
	query2 := "FROM 2021-12-22T15:33:13.0000005Z TO 2024-01-12T15:33:13.0000006Z  MATCH (a)-[x]->(b) WHERE b.properties_Risc > 0 RETURN  b, b.properties_Risc"
	query3 := "FROM 2021-12-22T15:33:13.0000005Z TO 2024-01-12T15:33:13.0000006Z  MATCH (a)-[x]->(b) WHERE a.properties_components_cpu = 'UGWJn' RETURN  a, a.properties_components_cpu"
	query4 := "FROM 2021-12-22T15:33:13.0000005Z TO 2024-01-12T15:33:13.0000006Z  MATCH (a)-[x]->(b) RETURN *"
	query5 := "FROM 2021-12-22T15:33:13.0000005Z TO 2024-01-12T15:33:13.0000006Z  MATCH (a)-[x]->(b) RETURN a.properties_components_cpu"
	query6 := "FROM 2021-12-22T15:33:13.0000005Z TO 2024-01-12T15:33:13.0000006Z  MATCH (a)-[x]->(b) WHERE a.properties_components_cpu = 'UGWJn' RETURN a, b"
	query7 := "FROM 2021-12-22T15:33:13.0000005Z TO 2024-01-12T15:33:13.0000006Z  MATCH (a)-[x]->(b) WHERE a.properties_components_cpu = 'TTT' RETURN *"
	query8 := "FROM 2021-12-22T15:33:13.0000005Z TO 2024-01-12T15:33:13.0000006Z  MATCH (a)-[x]->(b) WHERE b.properties_Risc > 40 RETURN *"
	query9 := "FROM 2021-12-22T15:33:13.0000005Z TO 2024-01-12T15:33:13.0000006Z  MATCH (a)-[x]->(b) WHERE b.properties_Risc > 39 RETURN *"

	expecteds := []string{expected1, expected2, expected3, expected4, expected5, expected6, expected7, expected8, expected9}
	queries := []string{query1, query2, query3, query4, query5, query6, query7, query8, query9}
	keys := [][]string{{"a", "x", "b"}, {"b", "b.properties_Risc"}, {"a", "a.properties_components_cpu"}, {"a", "b", "x"}, {"a.properties_components_cpu"}, {"a", "b"}, {}, {}, {"a", "b", "x"}}

	for i, query := range queries {
		queryInfo, err := parser.ParseQuery(cleanQuery(query))
		if err != nil {
			t.Fatalf("Error while parsing query: %v", err)
		}
		res, err := ProcessQuery(queryInfo)
		if err != nil {
			t.Fatalf("Error while processing query: %v", err)
		}

		// clean the internal neo4j IDs of the graph elements because the change on restore
		removeElementIDs(res)
		jsonRes := utils.JsonStringFromMapOrdered(res, keys[i])

		byteExpected := []byte(expecteds[i])
		bufferEx := new(bytes.Buffer)
		if err := json.Compact(bufferEx, byteExpected); err != nil {
			fmt.Println(err)
		}

		byteRes := []byte(jsonRes)
		bufferRes := new(bytes.Buffer)
		if err := json.Compact(bufferRes, byteRes); err != nil {
			fmt.Println(err)
		}

		if bufferEx.String() != bufferRes.String() {
			t.Fatalf("\nQuery: %v\nExpected\n  %v\nGot:\n  %v", query, bufferEx.String(), bufferRes.String())
		}
	}
}

func TestShallowQueries(t *testing.T) {

	// initialize Neo4j
	DriverNeo, err := neo4j.NewDriverWithContext("neo4j://"+testConfNeo.Host+":"+testConfNeo.Port, neo4j.BasicAuth(testConfNeo.Username, testConfNeo.Password, ""))
	if err != nil {
		log.Printf("Creating driver failed: %v", err)
		os.Exit(1)
	}
	defer DriverNeo.Close(context.Background())

	databaseapi.SessionNeo = DriverNeo.NewSession(context.Background(), neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer databaseapi.SessionNeo.Close(context.Background())

	// initialize TimescaleDB
	databaseapi.SessionTS, err = databaseapi.ConnectTimescale(testConfTS.Username, testConfTS.Password, testConfTS.Port, testConfTS.Database)
	if err != nil {
		log.Printf("Creating driver failed: %v", err)
		os.Exit(1)
	}

	defer databaseapi.SessionTS.Close()

	query1 := "FROM 2021-12-22T15:33:13.0000005Z TO 2024-01-12T15:33:13.0000006Z SHALLOW MATCH (a)-[x]->(b) WHERE a.properties_components_cpu IS NOT NULL RETURN *"
	query2 := "FROM 2021-12-22T15:33:13.0000005Z TO 2024-01-12T15:33:13.0000006Z SHALLOW MATCH (a)-[x]->(b) WHERE a.properties_components_cpu IS NOT NULL RETURN  a.properties_components_cpu"
	query3 := "FROM 2021-12-22T15:33:13.0000005Z TO 2024-01-12T15:33:13.0000006Z SHALLOW MATCH (a)-[x]->(b) WHERE a.notExistingProperty IS NOT NULL RETURN *"
	query4 := "FROM 2021-12-22T15:33:13.0000005Z TO 2024-01-12T15:33:13.0000006Z SHALLOW MATCH (a)-[x]->(b) WHERE b.properties_Risc > 0 RETURN  b, b.properties_Risc"
	query5 := "FROM 2021-12-22T15:33:13.0000005Z TO 2024-01-12T15:33:13.0000006Z SHALLOW MATCH (a)-[x]->(b) WHERE any(b.properties_Risc > 0) IS NOT NULL RETURN b, b.properties_Risc"
	query6 := "FROM 2021-12-22T15:33:13.0000005Z TO 2024-01-12T15:33:13.0000006Z SHALLOW MATCH (a)-[x]->(b) RETURN  a,x,b"
	query7 := "FROM 2021-12-22T15:33:13.0000005Z TO 2024-01-12T15:33:13.0000006Z SHALLOW MATCH (a)-[x]->(b) WHERE b.properties_Risc > 0 RETURN  b.properties_Risc"
	query8 := "FROM 2021-12-22T15:33:13.0000005Z TO 2024-01-12T15:33:13.0000006Z SHALLOW MATCH (a)-[x]->(b) WHERE a.properties_components_cpu = 'UGWJn' RETURN  a, a.properties_components_cpu"
	query9 := "FROM 2021-12-22T15:33:13.0000005Z TO 2024-01-12T15:33:13.0000006Z SHALLOW MATCH (a)-[x]->(b) RETURN  *"
	query10 := "FROM 2021-12-22T15:33:13.0000005Z TO 2024-01-12T15:33:13.0000006Z SHALLOW MATCH (a)-[x]->(b) RETURN  a.properties_components_cpu"
	query11 := "FROM 2021-12-22T15:33:13.0000005Z TO 2024-01-12T15:33:13.0000006Z  MATCH (a)-[x]->(b) WHERE a.properties_components_cpu = 'TTT' RETURN *"
	query12 := "FROM 2021-12-22T15:33:13.0000005Z TO 2024-01-12T15:33:13.0000006Z  MATCH (a)-[x]->(b) WHERE b.properties_Risc > 40 RETURN *"
	query13 := "FROM 2021-12-22T15:33:13.0000005Z TO 2024-01-12T15:33:13.0000006Z  MATCH (a)-[x]->(b) WHERE b.properties_Risc > 39 RETURN *"

	expecteds := []string{expectedShallow1, expectedShallow2, expectedShallow3, expectedShallow4, expectedShallow5, expectedShallow6, expectedShallow7, expectedShallow8, expectedShallow9, expectedShallow10, expectedShallow11, expectedShallow12, expectedShallow13}
	queries := []string{query1, query2, query3, query4, query5, query6, query7, query8, query9, query10, query11, query12, query13}
	keys := [][]string{{"a", "b", "x"}, {"a.properties_components_cpu"}, {"a", "b", "x"}, {"b", "b.properties_Risc"}, {"b", "b.properties_Risc"}, {"a", "x", "b"}, {"b.properties_Risc"}, {"a", "a.properties_components_cpu"}, {"a", "b", "x"}, {"a.properties_components_cpu"}, {}, {}, {"a", "b", "x"}}

	for i, query := range queries {
		queryInfo, err := parser.ParseQuery(cleanQuery(query))
		if err != nil {
			t.Fatalf("Error while parsing query: %v", err)
		}
		res, err := ProcessQuery(queryInfo)
		if err != nil {
			t.Fatalf("Error while processing query: %v", err)
		}

		// clean the internal neo4j IDs of the graph elements because the change on restore
		removeElementIDs(res)
		jsonRes := utils.JsonStringFromMapOrdered(res, keys[i])

		byteExpected := []byte(expecteds[i])
		bufferEx := new(bytes.Buffer)
		if err := json.Compact(bufferEx, byteExpected); err != nil {
			fmt.Println(err)
		}

		byteRes := []byte(jsonRes)
		bufferRes := new(bytes.Buffer)
		if err := json.Compact(bufferRes, byteRes); err != nil {
			fmt.Println(err)
		}

		if bufferEx.String() != bufferRes.String() {
			t.Fatalf("\nQuery: %v\nExpected\n  %v\nGot:\n  %v", query, bufferEx.String(), bufferRes.String())
		}
	}
}

func removeElementIDs(graphData interface{}) {
	switch value := graphData.(type) {
	default:
	case map[string]interface{}:
		for key, val := range value {
			if key == "ElementId" {
				delete(value, key) // Remove the "ElementId" field from the map
			} else {
				removeElementIDs(val)
			}
		}
	case map[string][]interface{}:
		for _, val := range value {
			removeElementIDs(val)
		}
	case []interface{}:
		for i, val := range value {
			removeElementIDs(val)
			switch el := val.(type) {
			default:
			case neo4j.Node:
				el.ElementId = ""
				value[i] = el
			case neo4j.Relationship:
				el.ElementId = ""
				el.EndElementId = ""
				el.StartElementId = ""
				value[i] = el
			}
		}
	}
}

// func printRes(t *testing.T, queryRes map[string][]any, queryInfo parser.ParseResult) {
// 	t.Logf("\n\n\n                 		 QUERY RESULT\n						%+v\n\n\n", queryRes)
// 	if len(queryInfo.ReturnProjections) > 0 {
// 		t.Log("\n\n\n                      Printed ordered                         \n\n\n\n")
// 		t.Logf("%+v\n", utils.JsonStringFromMapOrdered(queryRes, queryInfo.ReturnProjections))
// 	} else {
// 		t.Log("\n\n\n                      Printed unordered                         \n\n\n\n")
// 		t.Logf("%+v\n", utils.JsonStringFromMap(queryRes))
// 	}
// }
//
