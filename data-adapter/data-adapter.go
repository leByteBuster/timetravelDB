package dataadapter

import (
	"context"
	"fmt"
	"log"

	databaseapi "github.com/LexaTRex/timetravelDB/database-api"
	"github.com/LexaTRex/timetravelDB/utils"
	"github.com/google/uuid"
)

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
		// TODO: remove
		fmt.Println("Load template: ", template)
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

	utils.Debugf("\n Number nodes: %v", len(graph_nodes))
	utils.Debugf("\n Number edges: %v", len(graph_edges))

	nodeQueries, timeSeriesNodes := getQueryStringsNodes(graph_nodes)
	edgeQueries, timeSeriesEdges := getQueryStringsEdges(graph_edges)
	utils.Debugf("\n edge time-series: %v", timeSeriesEdges)

	databaseapi.WriteQueryMultipleNeo4j(context.Background(), nodeQueries, map[string]interface{}{})
	databaseapi.WriteQueryMultipleNeo4j(context.Background(), edgeQueries, map[string]interface{}{})

	utils.Debugf("\n node Queries: %v", nodeQueries)
	utils.Debugf("\n edge Queries: %v", edgeQueries)
	utils.Debugf("\n node time-series: %v", timeSeriesNodes)
	utils.Debugf("\n Number time-series map nodes: %v", len(timeSeriesNodes))
	utils.Debugf("\n Number time-series map edges: %v", len(timeSeriesEdges))

	timeSeries := map[uuid.UUID][]map[string]interface{}{}

	for k, v := range timeSeriesNodes {
		timeSeries[k] = v
	}
	for k, v := range timeSeriesEdges {
		timeSeries[k] = v
	}

	loadDataTimeScaleDB(timeSeries)
}
