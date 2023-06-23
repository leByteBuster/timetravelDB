package dataadapterneo4j

// this adapter allows to read in temporal property graph data in the form of json files generated by the data_generator as nodes
// and into a neo4j database

import (
	"context"
	"fmt"
	"log"
	"time"

	databaseapi "github.com/LexaTRex/timetravelDB/database-api"
	"github.com/LexaTRex/timetravelDB/utils"
)

type TmpPropVal[T any] struct {
	Start string
	End   string
	Value T
}

func LoadData(template string) {

	var graph_nodes []map[string]interface{}
	var graph_edges []map[string]interface{}
	var err error

	if template == "" {
		graph_nodes, err = utils.LoadJsonData("data-generator/generated-data/graph_nodes.json")
		if err != nil {
			log.Printf("Error loading nodes from json: %v", err)
		}
		graph_edges, err = utils.LoadJsonData("data-generator/generated-data/graph_edges.json")
		if err != nil {
			log.Printf("Error loading edges from json: %v", err)
			return
		}
	} else {
		graph_nodes, err = utils.LoadJsonData("data-generator/generated-data/graph_nodes" + template + ".json")
		if err != nil {
			log.Printf("Error loading nodes from json: %v", err)
		}
		graph_edges, err = utils.LoadJsonData("data-generator/generated-data/graph_edges" + template + ".json")
		if err != nil {
			log.Printf("Error loading edges from json: %v", err)
			return
		}
	}

	loadGraphNodesIntoNeo4jDatabase(graph_nodes, context.Background())
	loadGraphEdgesIntoNeo4jDatabase(graph_edges, context.Background())

	if err != nil {
		log.Printf("Error: %v", err)
	}
}

func loadGraphNodesIntoNeo4jDatabase(graph_nodes []map[string]interface{}, ctx context.Context) {

	// Loop through the data array
	for _, node := range graph_nodes {

		labels := node["labels"]

		// temporality
		start := node["start"]
		end := node["end"]

		delete(node, "labels")
		delete(node, "start")
		delete(node, "end")

		queryFlat := `CREATE (n:` + labels.([]interface{})[0].(string) + ` {start: datetime("` + start.(string) + `"), end: datetime("` + end.(string) + `"),`
		primaryQueryFragmentsFlat := generateNeo4jFlatProperties(node)

		for _, fragment := range primaryQueryFragmentsFlat {
			queryFlat += fragment
		}

		queryFlat = queryFlat[:len(queryFlat)-2] + `})`

		utils.Debugf("\nNode Query: %v\n", queryFlat)

		databaseapi.WriteQueryNeo4j(ctx, queryFlat, map[string]interface{}{})
	}
}

func loadGraphEdgesIntoNeo4jDatabase(graph_edges []map[string]interface{}, ctx context.Context) {

	// Loop through the data array
	for _, edge := range graph_edges {

		// temporality
		start := edge["start"]
		end := edge["end"]

		// vector
		from := edge["from"]
		to := edge["to"]

		label := edge["label"]

		delete(edge, "from")
		delete(edge, "to")
		delete(edge, "start")
		delete(edge, "end")
		delete(edge, "labels")

		queryPrefix := `MATCH (a),(b) WHERE a.nodeid = $from AND b.nodeid = $to CREATE (a)-[r:` + label.(string) + ` {start: datetime("` + start.(string) + `"), end: datetime("` + end.(string) + `"),`
		querySuffix := `}]->(b)`

		neo4jEdgeProperties := generateNeo4jFlatProperties(edge)

		queryProperties := ""

		for _, fragment := range neo4jEdgeProperties {
			queryProperties += fragment
		}

		queryProperties = queryProperties[:len(queryProperties)-2]

		edgeQuery := queryPrefix + queryProperties + querySuffix

		utils.Debugf("\nEdge Query: %v\n", edgeQuery)
		utils.Debugf("\nFrom: %v\n", from)
		utils.Debugf("\nTo: %v\n", to)

		databaseapi.WriteQueryNeo4j(ctx, edgeQuery, map[string]interface{}{"from": from, "to": to})
	}
}

func generateNeo4jFlatProperties(property map[string]interface{}) []string {
	queryBaseFragment := ""
	queryPropertyFragments := make([]string, 0)
	for key, value := range property {
		switch propertyValue := value.(type) {

		case map[string]interface{}:
			propertyFragments := generateNeo4jFlatProperties(propertyValue)
			for _, fragment := range propertyFragments {
				queryPropertyFragments = append(queryPropertyFragments, queryBaseFragment+key+`_`+fragment)
			}

		// in the case of string not much has to be done but this case does not occour because the values are always lists of maps/objects
		case string:

			propertyFragment := key + `: "` + propertyValue + `", `
			queryPropertyFragments = append(queryPropertyFragments, propertyFragment)

		// in the case of an array of maps, every map object represents a temporal value of a property in the form of {start:_, end:_, value:_}. Create an
		// query entry for creating an own node for every of these temporal property values
		case []interface{}:

			valueList := utils.ConvertMaps(propertyValue)

			// iterate over the array of maps and create a value-node query fragment for each of them (for the CREATE query)
			for i, value := range valueList {

				_, err := time.Parse("2006-01-02T15:04:05.9999999999Z", fmt.Sprint(value["Start"]))
				_, err2 := time.Parse("2006-01-02T15:04:05.9999999999Z", fmt.Sprint(value["End"]))

				if err != nil || err2 != nil {
					panic("\nError Parsing DateTime. No ISO8601")
				}

				// generate a unique property entry for every property value in the list. Number the property fields by the index of the list
				valueFragmentTime := key + `_` + fmt.Sprint(i) + `_` + `time` + `: datetime("` + fmt.Sprint(value["Start"]) + `"), `
				// valueFragmentStart := key + `_` + fmt.Sprint(i) + `_` + `start` + `: "` + fmt.Sprint(value["Start"]) + `", `
				// valueFragmentEnd := key + `_` + fmt.Sprint(i) + `_` + `end` + `: "` + fmt.Sprint(value["End"]) + `", `
				valueFragmentValue := key + `_` + fmt.Sprint(i) + `_` + `value` + `: `
				switch valueType := value["Value"].(type) {
				case string:
					valueFragmentValue += `"` + valueType + `", `
				default:
					valueFragmentValue += fmt.Sprint(valueType) + `, `
				}

				//queryPropertyFragments = append(queryPropertyFragments, []string{valueFragmentStart, valueFragmentEnd, valueFragmentValue}...)
				queryPropertyFragments = append(queryPropertyFragments, []string{valueFragmentTime, valueFragmentValue}...)
			}

		default:
			propertyFragment := key + `: ` + fmt.Sprint(propertyValue) + `, `
			queryPropertyFragments = append(queryPropertyFragments, propertyFragment)
		}
	}
	return queryPropertyFragments
}
