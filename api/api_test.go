package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/LexaTRex/timetravelDB/utils"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

/**
  * This is a test integration test file.
  * It tests the API by sending requests to the API and checking the response.
  * For the tests to work, there must be a Neo4j database running on the default port,
  * as well as a timetravelDB database running on the default port.
  * See DOCKER_README.md for setting up the testing environment.
**/

func TestNonShallowQueries(t *testing.T) {

	PassNeo = "test"
	UserNeo = "neo4j"

	UserTS = "postgres"
	PassTS = "password"
	DBnameTS = "postgres"

	var err error

	// initialize Neo4j
	DriverNeo, err = neo4j.NewDriverWithContext(UriNeo, neo4j.BasicAuth(UserNeo, PassNeo, ""))
	if err != nil {
		log.Printf("Creating driver failed: %v", err)
		os.Exit(1)
	}
	defer DriverNeo.Close(context.Background())

	SessionNeo = DriverNeo.NewSession(context.Background(), neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer SessionNeo.Close(context.Background())

	// initialize TimescaleDB
	SessionTS, err = connectTimescale(UserTS, PassTS, PortTS, DBnameTS)
	if err != nil {
		log.Printf("Creating driver failed: %v", err)
		os.Exit(1)
	}
	defer SessionTS.Close(context.Background())

	query1 := "FROM 2021-12-22T15:33:13.0000005Z TO 2024-01-12T15:33:13.0000006Z  MATCH (a)-[x]->(b) RETURN  a,x,b"
	query2 := "FROM 2021-12-22T15:33:13.0000005Z TO 2024-01-12T15:33:13.0000006Z  MATCH (a)-[x]->(b) WHERE b.properties_Risc > 0 RETURN  b, b.properties_Risc"
	query3 := " FROM 2021-12-22T15:33:13.0000005Z TO 2024-01-12T15:33:13.0000006Z  MATCH (a)-[x]->(b) WHERE a.properties_components_cpu = 'UGWJn' RETURN  a, a.properties_components_cpu"
	query4 := "FROM 2021-12-22T15:33:13.0000005Z TO 2024-01-12T15:33:13.0000006Z  MATCH (a)-[x]->(b) RETURN  *"
	expecteds := []string{expected1, expected2, expected3, expected4}
	queries := []string{query1, query2, query3, query4}
	keys := [][]string{{"a", "x", "b"}, {"b", "b.properties_Risc"}, {"a", "a.properties_components_cpu"}, {"a", "b", "x"}}

	for i, query := range queries {
		res, err := ProcessQuery(cleanQuery(query))
		fmt.Printf("res before removing: %+v", res)
		removeElementIDs(res)
		jsonRes := utils.JsonStringFromMapOrdered(res, keys[i])
		if err != nil {
			t.Fatalf("Error while processing query: %v", err)
		}
		if err != nil {
			t.Fatalf("Error while processing query: %v", err)
		}
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

// TODO
func TestShallowQueries(t *testing.T) {

	PassNeo = "test"
	UserNeo = "neo4j"

	UserTS = "postgres"
	PassTS = "password"
	DBnameTS = "postgres"

	var err error

	// initialize Neo4j
	DriverNeo, err = neo4j.NewDriverWithContext(UriNeo, neo4j.BasicAuth(UserNeo, PassNeo, ""))
	if err != nil {
		log.Printf("Creating driver failed: %v", err)
		os.Exit(1)
	}
	defer DriverNeo.Close(context.Background())

	SessionNeo = DriverNeo.NewSession(context.Background(), neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer SessionNeo.Close(context.Background())

	// initialize TimescaleDB
	SessionTS, err = connectTimescale(UserTS, PassTS, PortTS, DBnameTS)
	if err != nil {
		log.Printf("Creating driver failed: %v", err)
		os.Exit(1)
	}
	defer SessionTS.Close(context.Background())

	query1 := "FROM 2021-12-22T15:33:13.0000005Z TO 2024-01-12T15:33:13.0000006Z SHALLOW MATCH (a)-[x]->(b) WHERE a.properties_components_cpu IS NOT NULL RETURN *"
	query2 := "FROM 2021-12-22T15:33:13.0000005Z TO 2024-01-12T15:33:13.0000006Z SHALLOW MATCH (a)-[x]->(b) WHERE a.properties_components_cpu IS NOT NULL RETURN  a.properties_components_cpu"
	query3 := "FROM 2021-12-22T15:33:13.0000005Z TO 2024-01-12T15:33:13.0000006Z SHALLOW MATCH (a)-[x]->(b) WHERE a.notExistingProperty IS NOT NULL RETURN *"
	query4 := "FROM 2021-12-22T15:33:13.0000005Z TO 2024-01-12T15:33:13.0000006Z SHALLOW MATCH (a)-[x]->(b) WHERE b.properties_Risc > 0 RETURN  b, b.properties_Risc"
	query5 := "FROM 2021-12-22T15:33:13.0000005Z TO 2024-01-12T15:33:13.0000006Z SHALLOW MATCH (a)-[x]->(b) WHERE any(b.properties_Risc > 0) IS NOT NULL RETURN b, b.properties_Risc"
	query6 := "FROM 2021-12-22T15:33:13.0000005Z TO 2024-01-12T15:33:13.0000006Z SHALLOW MATCH (a)-[x]->(b) RETURN  a,x,b"
	query7 := "FROM 2021-12-22T15:33:13.0000005Z TO 2024-01-12T15:33:13.0000006Z SHALLOW MATCH (a)-[x]->(b) WHERE b.properties_Risc > 0 RETURN  b.properties_Risc"
	query8 := "FROM 2021-12-22T15:33:13.0000005Z TO 2024-01-12T15:33:13.0000006Z SHALLOW MATCH (a)-[x]->(b) WHERE a.properties_components_cpu = 'UGWJn' RETURN  a, a.properties_components_cpu"
	query9 := "FROM 2021-12-22T15:33:13.0000005Z TO 2024-01-12T15:33:13.0000006Z SHALLOW MATCH (a)-[x]->(b) RETURN  *"
	expecteds := []string{expectedShallow1, expectedShallow2, expectedShallow3, expectedShallow4, expectedShallow5, expectedShallow6, expectedShallow7, expectedShallow8, expectedShallow9}
	queries := []string{query1, query2, query3, query4, query5, query6, query7, query8, query9}
	keys := [][]string{{"a", "b", "x"}, {"a.properties_components_cpu"}, {"a", "b", "x"}, {"b", "b.properties_Risc"}, {"b", "b.properties_Risc"}, {"a", "x", "b"}, {"b.properties_Risc"}, {"a", "a.properties_components_cpu"}, {"a", "b", "x"}}

	for i, query := range queries {
		res, err := ProcessQuery(cleanQuery(query))
		removeElementIDs(res)
		jsonRes := utils.JsonStringFromMapOrdered(res, keys[i])
		if err != nil {
			t.Fatalf("Error while processing query: %v", err)
		}

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
