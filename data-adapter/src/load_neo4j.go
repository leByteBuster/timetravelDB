package main

import (
	"fmt"

	"github.com/google/uuid"
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

func getQuerStringsNodes(graph_nodes []map[string]interface{}) ([]string, map[uuid.UUID][]map[string]interface{}) {

	var queries = make([]string, 0)
	var timeSeries = map[uuid.UUID][]map[string]interface{}{}

	// Loop through the data array
	for _, node := range graph_nodes {

		// TODO: node["label"] should be able to contain a list of labels
		query := `CREATE (n:` + node["label"].(string) + ` {`

		delete(node, "label")

		primaryQueryFragmentsFlat, timeSeriesMapNode := generateNeo4jFlatProperties(node)

		for k, v := range timeSeriesMapNode {
			timeSeries[k] = v
		}

		var propertyQueryString = ""

		for _, fragment := range primaryQueryFragmentsFlat {
			propertyQueryString += fragment
		}

		query += propertyQueryString

		query = query[:len(query)-2] + `})`

		queries = append(queries, query)
	}
	return queries, timeSeries
}

func getQueryStringsEdges(graph_edges []map[string]interface{}) ([]string, map[uuid.UUID][]map[string]interface{}) {

	var queries = make([]string, 0)
	var timeSeries = map[uuid.UUID][]map[string]interface{}{}

	// Loop through the data array
	for _, edge := range graph_edges {

		from := edge["from"]
		to := edge["to"]

		// maybe make a copy forst to keep the original ?  or is the map copied into this funciton anyways ?
		delete(edge, "from")
		delete(edge, "to")

		queryPrefix := `MATCH (a),(b) WHERE a.nodeid = ` + fmt.Sprint(from) + ` AND b.nodeid = ` + fmt.Sprint(to) + ` CREATE (a)-[r:Relation {`
		querySuffix := `}]->(b)`

		neo4jEdgeProperties, timeSeriesMapEdge := generateNeo4jFlatProperties(edge)

		for k, v := range timeSeriesMapEdge {
			timeSeries[k] = v
		}

		queryProperties := ""

		for _, fragment := range neo4jEdgeProperties {
			queryProperties += fragment
		}

		queryProperties = queryProperties[:len(queryProperties)-2]

		edgeQuery := queryPrefix + queryProperties + querySuffix
		queries = append(queries, edgeQuery)

	}

	return queries, timeSeries
}

func generateNeo4jFlatProperties(property map[string]interface{}) ([]string, map[uuid.UUID][]map[string]interface{}) {

	queryBaseFragment := ""
	queryPropertyFragments := make([]string, 0)

	timeSeriesMap := map[uuid.UUID][]map[string]interface{}{}
	for key, value := range property {
		switch propertyValue := value.(type) {

		case map[string]interface{}:
			propertyFragments, timeSeriesMapTmp := generateNeo4jFlatProperties(propertyValue)
			for k, v := range timeSeriesMapTmp {
				timeSeriesMap[k] = v
			}
			for _, fragment := range propertyFragments {
				queryPropertyFragments = append(queryPropertyFragments, queryBaseFragment+key+`_`+fragment)
			}

		case []interface{}:
			valueList := convertMaps(propertyValue)

			id := uuid.New()
			valueFragmentUuid := key + `: "` + id.String() + `", `
			queryPropertyFragments = append(queryPropertyFragments, valueFragmentUuid)
			timeSeriesMap[id] = valueList
		case string:
			propertyFragment := key + `: "` + propertyValue + `", `
			queryPropertyFragments = append(queryPropertyFragments, propertyFragment)
		case interface{}:
			propertyFragment := key + `: ` + fmt.Sprint(propertyValue) + `, `
			queryPropertyFragments = append(queryPropertyFragments, propertyFragment)
		default:
			panic("should not happen")
		}
	}
	return queryPropertyFragments, timeSeriesMap
}
