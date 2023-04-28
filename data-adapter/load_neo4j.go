package dataadapter

import (
	"fmt"

	"github.com/LexaTRex/timetravelDB/utils"
	"github.com/google/uuid"
)

func getQueryStringsNodes(graph_nodes []map[string]interface{}) ([]string, map[uuid.UUID][]map[string]interface{}) {

	var queries = make([]string, 0)
	var timeSeries = map[uuid.UUID][]map[string]interface{}{}

	// Loop through the data array
	for _, node := range graph_nodes {

		// temporality
		start := node["start"]
		end := node["end"]

		labels := node["labels"]

		delete(node, "start")
		delete(node, "end")

		delete(node, "labels")

		// TODO: node["label"] should be able to contain a list of labels
		query := `CREATE (n:` + labels.([]interface{})[0].(string) + ` {`

		queryTemporalProperties := `start: datetime('` + fmt.Sprint(start) + `'), end: datetime('` + fmt.Sprint(end) + `'),`

		primaryQueryFragmentsFlat, timeSeriesMapNode := generateNeo4jFlatProperties(node)

		for k, v := range timeSeriesMapNode {
			timeSeries[k] = v
		}

		var propertyQueryString = ""

		for _, fragment := range primaryQueryFragmentsFlat {
			propertyQueryString += fragment
		}

		query += queryTemporalProperties + propertyQueryString

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

		// temporality
		start := edge["start"]
		end := edge["end"]

		// vector
		from := edge["from"]
		to := edge["to"]

		label := edge["label"]

		// maybe make a copy forst to keep the original ?  or is the map copied into this funciton anyways ?

		delete(edge, "start")
		delete(edge, "end")

		delete(edge, "from")
		delete(edge, "to")

		delete(edge, "label")

		queryPrefix := `MATCH (a),(b) WHERE a.nodeid = ` + fmt.Sprint(from) + ` AND b.nodeid = ` + fmt.Sprint(to) + ` CREATE (a)-[r:` + label.(string) + ` {`
		queryTemporalProperties := `start: datetime('` + fmt.Sprint(start) + `'), end: datetime('` + fmt.Sprint(end) + `'),`
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

		edgeQuery := queryPrefix + queryTemporalProperties + queryProperties + querySuffix

		queries = append(queries, edgeQuery)

		utils.Debugf("Edge Query: %v", edgeQuery)
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
			valueList := utils.ConvertMaps(propertyValue)

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
