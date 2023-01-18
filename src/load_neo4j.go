package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type comp struct {
	id        int16
	some_data string
	is_ts     bool
}

type TmpPropVal[T any] struct {
	Start string
	End   string
	Value T
}

func loadGraphNodesIntoNeo4jDatabase(graph_nodes []map[string]interface{}, ctx context.Context, uri, username, password string) map[uuid.UUID][]map[string]interface{} {

	var timeSeries = map[uuid.UUID][]map[string]interface{}{}

	// Connect to the Neo4j database
	driver, err := neo4j.NewDriverWithContext(uri, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer driver.Close(ctx)

	session := driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	// Loop through the data array
	for _, node := range graph_nodes {

		// Build the CREATE query
		// primaryQueryNode := `CREATE (:Node {`
		queryFlat := `CREATE (:Node {`

		primaryQueryFragmentsFlat, timeSeriesMapNode := generateNeo4jFlatProperties(node)
		for k, v := range timeSeriesMapNode {
			timeSeries[k] = v
		}

		var propertyQueryString = ""

		for _, fragment := range primaryQueryFragmentsFlat {
			propertyQueryString += fragment
		}

		queryFlat += propertyQueryString

		queryFlat = queryFlat[:len(queryFlat)-2] + `})`

		fmt.Printf("\nFlat Query: %v\n", queryFlat)
		fmt.Printf("\nTime-Series Map: %v\n", timeSeriesMapNode)

		res, err := session.Run(ctx, queryFlat, map[string]interface{}{})
		if err != nil {
			fmt.Println(err)
			return nil
		}
		fmt.Println("Result: ", res)
	}
	return timeSeries
}

func loadGraphEdgesIntoNeo4jDatabase(graph_edges []map[string]interface{}, ctx context.Context, uri, username, password string) map[uuid.UUID][]map[string]interface{} {

	var timeSeries = map[uuid.UUID][]map[string]interface{}{}

	// Connect to the Neo4j database
	driver, err := neo4j.NewDriverWithContext(uri, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer driver.Close(ctx)

	session := driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	// Loop through the data array
	for _, edge := range graph_edges {

		from := edge["from"]
		to := edge["to"]

		// maybe make a copy forst to keep the original ?  or is the map copied into this funciton anyways ?
		delete(edge, "from")
		delete(edge, "to")

		queryPrefix := `MATCH (a),(b) WHERE a.nodeid = $from AND b.nodeid = $to CREATE (a)-[r:Relation {`
		querySuffix := `}]->(b)`

		neo4jEdgeProperties, timeSeriesMapNode := generateNeo4jFlatProperties(edge)

		for k, v := range timeSeriesMapNode {
			timeSeries[k] = v
		}

		queryProperties := ""

		for _, fragment := range neo4jEdgeProperties {
			queryProperties += fragment
		}

		queryProperties = queryProperties[:len(queryProperties)-2]

		edgeQuery := queryPrefix + queryProperties + querySuffix

		fmt.Printf("\nEdge Query: %v\n", edgeQuery)
		fmt.Printf("\nFrom: %v\n", from)
		fmt.Printf("\nTo: %v\n", to)
		fmt.Printf("\nTimeSeries Map: %v\n", timeSeriesMapNode)

		// res, err := session.Run(ctx, ``,
		//	map[string]interface{}{})
		// Execute the CREATE query
		res, err := session.Run(ctx, edgeQuery, map[string]interface{}{"from": from, "to": to})
		if err != nil {
			fmt.Println(err)
			return nil
		}
		fmt.Println("Result: ", res)
	}

	return timeSeries
}

func generateNeo4jFlatProperties(property map[string]interface{}) ([]string, map[uuid.UUID][]map[string]interface{}) {
	queryBaseFragment := ""
	queryPropertyFragments := make([]string, 0)
	timeSeriesMap := map[uuid.UUID][]map[string]interface{}{}
	for key, value := range property {
		switch propertyValue := value.(type) {

		case map[string]interface{}:
			propertyFragments, timeSeriesMapTmp := generateNeo4jFlatProperties(propertyValue)
			timeSeriesMap = timeSeriesMapTmp
			for _, fragment := range propertyFragments {
				queryPropertyFragments = append(queryPropertyFragments, queryBaseFragment+key+`_`+fragment)
			}

		case []interface{}:

			valueList := convertMaps(propertyValue)

			id := uuid.New()

			// generate a unique property entry for every property value in the list. Number the property fields by the index of the list
			valueFragmentUuid := key + `: "` + id.String() + `", `

			queryPropertyFragments = append(queryPropertyFragments, valueFragmentUuid)

			// map uiid to valuelist
			timeSeriesMap[id] = valueList

		// the next two cases are only supposed to happend if the input data is already formated so the values are references
		// to time-series lists in the timescaledb
		// in the case of string not much has to be done but this case does not occour because the values are always lists of maps/objects
		case string:
			propertyFragment := key + `: "` + propertyValue + `", `
			queryPropertyFragments = append(queryPropertyFragments, propertyFragment)

		// in the case of an array of maps, every map object represents a temporal value of a property in the form of {start:_, end:_, value:_}. Create an
		// query entry for creating an own node for every of these temporal property values

		// in the case of interface not much has to be done but this case does not occour because the values are always lists of maps/objects
		case interface{}:
			propertyFragment := key + `: ` + fmt.Sprint(propertyValue) + `, `
			queryPropertyFragments = append(queryPropertyFragments, propertyFragment)

		// in the case of an array of maps, every map object represents a temporal value of a property in the form of {start:_, end:_, value:_}. Create an
		// query entry for creating an own node for every of these temporal property values

		default:
			panic("should not happen")
		}
	}
	return queryPropertyFragments, timeSeriesMap
}

func loadJsonData(path string) ([]map[string]interface{}, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Decode the JSON data
	var data []map[string]interface{}
	err = json.NewDecoder(file).Decode(&data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func convertMaps(originalMaps []interface{}) []map[string]interface{} {
	convertedMaps := make([]map[string]interface{}, 0)
	for _, originalMap := range originalMaps {
		convertedMap := map[string]interface{}{}
		for key, value := range originalMap.(map[string]interface{}) {
			convertedMap[key] = value.(interface{})
		}
		convertedMaps = append(convertedMaps, convertedMap)
	}
	return convertedMaps
}

// ### old function from testprojeckt ###
func load_data(ctx context.Context, uri, username, password string, cmp []comp) (string, error) {
	driver, err := neo4j.NewDriverWithContext(uri, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		return "", err
	}
	defer driver.Close(ctx)
	fmt.Println("check 2")

	session := driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	// fmt.Println("check 3")
	// _, err = session.ExecuteWrite(ctx, func(transaction neo4j.ManagedTransaction) (any, error) {
	// 	result, err := transaction.Run(ctx,
	// 		"CREATE (a:Computer) SET a.message = $message",
	// 		map[string]any{"message": "available"})
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	if result.Next(ctx) {
	// 		return result.Record().Values[0], nil
	// 	}
	// 	fmt.Println("check 4")

	// 	return nil, result.Err()
	// })

	// _, err = session.ExecuteWrite(ctx, func(transaction neo4j.ManagedTransaction) (any, error) {
	// 	fmt.Println("check 5.1")
	// 	result, err := transaction.Run(ctx, "WITH [{id:1,some_data: 'testdata',is_ts:true},{id:2,some_data:'testdata',is_ts:false}] as nodes UNWIND nodes AS node MERGE (computer:Computer {id: node.id, someData: node.some_data, isTs: node.is_ts }) SET computer += node RETURN computer",
	// 		map[string]any{})

	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	if result.Next(ctx) {
	// 		return result.Record().Values[0], nil
	// 	}
	// 	fmt.Println("check 4")

	// 	return nil, result.Err()
	// })

	_, err = session.ExecuteWrite(ctx, func(transaction neo4j.ManagedTransaction) (any, error) {
		fmt.Println("check 5.1")
		comp_list := []map[string]interface{}{
			{
				"id":        11,
				"some_data": "testdata",
				"is_ts":     true,
			},
			{
				"id":        22,
				"some_data": "testdata",
				"is_ts":     false,
			},
		}
		result, err := transaction.Run(ctx, "UNWIND $nodes AS node MERGE (computer:Computer {id: node.id, some_data: node.some_data, is_ts: node.is_ts }) SET computer += node RETURN computer",
			map[string]any{"nodes": comp_list})

		if err != nil {
			fmt.Println("check 5.2")
			fmt.Printf("Error: %v", err)
			return nil, err
		}

		if result.Next(ctx) {
			fmt.Println("check 5.3")
			return result.Record().Values[0], nil
		}
		fmt.Println("check 4")

		return nil, result.Err()
	})

	if err != nil {
		return "", err
	}

	fmt.Println("check 5")
	return "success", nil
}
